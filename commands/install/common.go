// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package install

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk/common"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	UserLocalBin     = "/usr/local/bin/"
	UserLocalInclude = "/usr/local/include/"
)

var (
	GithubClient = github.NewClient(nil)
	RkHomeDir    = "./"
	GithubInfo   = githubInfo{
		DecompressedFilesToCopy:      make([]string, 0),
		DecompressedFilesDestination: make([]string, 0),
	}
)

func init() {
	userHomeDir, _ := os.UserHomeDir()
	RkHomeDir = path.Join(userHomeDir, ".rk")
	os.MkdirAll(RkHomeDir, os.ModePerm)
}

type githubInfo struct {
	Owner                        string                      `json:"owner"`
	Repo                         string                      `json:"repo"`
	GoGetUrl                     string                      `json:"-"`
	ReleaseListFromGithub        []*github.RepositoryRelease `json:"-"`
	ReleaseToDownload            *github.RepositoryRelease   `json:"-"`
	DownloadUrl                  string                      `json:"downloadUrl"`
	DownloadTo                   string                      `json:"-"`
	DownloadFilter               func(string) bool           `json:"-"`
	DecompressTo                 string                      `json:"decompressTo"`
	DecompressType               string                      `json:"-"`
	DecompressedFilesToCopy      []string                    `json:"-"`
	DecompressedFilesDestination []string                    `json:"decompressedFilesDestination"`
	ValidationCmd                *exec.Cmd                   `json:"-"`
}

func commandDefault(name string) *cli.Command {
	return &cli.Command{
		Name:      name,
		Usage:     fmt.Sprintf("install %s on local machine", name),
		UsageText: fmt.Sprintf("rk install %s -r [release]", name),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "tag",
				Aliases: []string{"t"},
				Usage:   "specify which tag to install, will use latest one if not specified",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list releases, list most recent 10 releases",
			},
		},
	}
}

func beforeDefault(ctx *cli.Context) error {
	name := strings.Join(strings.Split(ctx.Command.FullName(), " "), "/")
	event := common.CreateEvent(name)
	event.AddPayloads(zap.Strings("flags", ctx.FlagNames()))

	// Inject event into context
	ctx.Context = context.WithValue(ctx.Context, common.EventKey, event)

	return nil
}

func afterDefault(ctx *cli.Context) error {
	common.Finish(common.GetEvent(ctx), nil)
	return nil
}

func hasListFlag(ctx *cli.Context) bool {
	return ctx.IsSet("list")
}

func getTagFromFlags(ctx *cli.Context) string {
	if ctx.IsSet("tag") {
		return ctx.Value("tag").(string)
	}

	return ""
}

func githubInfoToPayloads() []zap.Field {
	return []zap.Field{
		zap.String("owner", GithubInfo.Owner),
		zap.String("repo", GithubInfo.Repo),
		zap.String("url", GithubInfo.DownloadUrl),
		zap.String("tempPath", GithubInfo.DecompressTo),
		zap.Strings("destination", GithubInfo.DecompressedFilesDestination),
	}
}

func listTagsFromGithub(ctx *cli.Context) error {
	opt := &github.ListOptions{PerPage: 5}
	res, _, err := GithubClient.Repositories.ListReleases(
		context.Background(),
		GithubInfo.Owner,
		GithubInfo.Repo,
		opt)

	if err != nil {
		return err
	}

	GithubInfo.ReleaseListFromGithub = res

	return nil
}

func printTagsFromGithub(ctx *cli.Context) error {
	opt := &github.ListOptions{PerPage: 5}
	res, _, err := GithubClient.Repositories.ListReleases(
		context.Background(),
		GithubInfo.Owner,
		GithubInfo.Repo,
		opt)

	if err != nil {
		return err
	}

	if len(res) < 1 {
		color.Yellow("No tags found")
		return nil
	}

	for i := range res {
		color.Yellow("- %s  (%s)", res[i].GetTagName(), res[i].GetPublishedAt().String())
	}

	return nil
}

func getReleaseToInstallFromGithub(ctx *cli.Context) error {
	tagFromFlags := getTagFromFlags(ctx)

	// User provided tag to download, we need to get release info from github
	if len(tagFromFlags) > 0 {
		res, _, err := GithubClient.Repositories.GetReleaseByTag(
			context.Background(),
			GithubInfo.Owner,
			GithubInfo.Repo,
			tagFromFlags)

		if err != nil {
			return err
		}

		GithubInfo.ReleaseToDownload = res

		color.White("- %s (%s)",
			GithubInfo.ReleaseToDownload.GetTagName(),
			GithubInfo.ReleaseToDownload.GetPublishedAt().String())
		return nil
	}

	// List tags and pick the latest one
	if err := listTagsFromGithub(ctx); err != nil {
		return err
	}

	if len(GithubInfo.ReleaseListFromGithub) < 1 {
		return errors.New("no releases available from github")
	}

	GithubInfo.ReleaseToDownload = GithubInfo.ReleaseListFromGithub[0]

	color.White("- %s (%s)",
		GithubInfo.ReleaseToDownload.GetTagName(),
		GithubInfo.ReleaseToDownload.GetPublishedAt().String())
	return nil
}

func downloadFromGithub(ctx *cli.Context) error {
	var id = int64(-1)
	var fileSize int64
	// 1: iterate assets in release and pick the right one
	for i := range GithubInfo.ReleaseToDownload.Assets {
		asset := GithubInfo.ReleaseToDownload.Assets[i]
		if GithubInfo.DownloadFilter(*asset.BrowserDownloadURL) {
			id = asset.GetID()
			fileSize = int64(asset.GetSize())
			GithubInfo.DownloadUrl = *asset.BrowserDownloadURL
			break
		}
	}

	// 2: construct local path of where to write the file
	GithubInfo.DownloadTo = path.Join(RkHomeDir, path.Base(GithubInfo.DownloadUrl))
	if id == -1 {
		return errors.New(fmt.Sprintf("failed to get downloadable id with tag:%s", GithubInfo.ReleaseToDownload.GetName()))
	}

	// 3: get reader object from github
	color.White("- Start to download %s to %s", GithubInfo.DownloadUrl, GithubInfo.DownloadTo)
	reader, _, err := GithubClient.Repositories.DownloadReleaseAsset(
		context.Background(),
		GithubInfo.Owner,
		GithubInfo.Repo,
		id,
		http.DefaultClient)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to call github API. %v", err))
	}

	// 4: create and open/create local file based on download URL
	f, err := os.OpenFile(GithubInfo.DownloadTo, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errors.New(fmt.Sprintf("failed to open file. %v", err))
	}
	defer f.Close()

	// 5: start downloading
	bar := progressbar.DefaultBytes(fileSize, "downloading")
	_, err = io.Copy(io.MultiWriter(f, bar), reader)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to write to file. %v", err))
	}

	return nil
}

func decompressFile(ctx *cli.Context) error {
	// 1: create temp dir
	if err := createTempDirForDecompress(); err != nil {
		return err
	}
	color.White("- Create temporary directory at %s", GithubInfo.DecompressTo)

	// 2: decompress
	color.White("- Decompress %s to %s", GithubInfo.DownloadTo, GithubInfo.DecompressTo)
	var err error
	if GithubInfo.DecompressType == "tar" {
		_, err = exec.Command("tar", "-C", GithubInfo.DecompressTo, "-xvzf", GithubInfo.DownloadTo).CombinedOutput()
	} else if GithubInfo.DecompressType == "zip" {
		_, err = exec.Command("unzip", "-o", GithubInfo.DownloadTo, "-d", GithubInfo.DecompressTo).CombinedOutput()
	} else {
		return errors.New("unknown file type to decompress")
	}

	if err != nil {
		return errors.New(fmt.Sprintf("failed to decompress. %v", err))
	}

	return nil
}

func copyToDestination(ctx *cli.Context) error {
	for i := range GithubInfo.DecompressedFilesDestination {
		srcPath := path.Join(GithubInfo.DecompressTo, GithubInfo.DecompressedFilesToCopy[i])
		destPath := GithubInfo.DecompressedFilesDestination[i]

		if GithubInfo.Repo == "protoc" {
			_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("chmod -R +rx %s", srcPath)).CombinedOutput()
			if err != nil {
				return errors.New(fmt.Sprintf("failed to chmod protoc. %v", err))
			}
		}
		_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo cp -r %s %s", srcPath, destPath)).CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("failed to copy %s to %s. %v", GithubInfo.DecompressTo, destPath, err))
		}

		color.White("- Copy %s to %s success", srcPath, destPath)
	}
	//
	//// 1: copy to target
	//srcPath := path.Join(GithubInfo.DecompressTo, GithubInfo.DecompressedFilesToCopy)
	//if GithubInfo.Repo == "protoc" {
	//	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("chmod -R +rx %s", srcPath)).CombinedOutput()
	//	if err != nil {
	//		return errors.New(fmt.Sprintf("failed to chmod protoc. %v", err))
	//	}
	//}
	//_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo cp -r %s %s", srcPath, GithubInfo.DecompressedFilesDestination)).CombinedOutput()
	//if err != nil {
	//	return errors.New(fmt.Sprintf("failed to copy %s to %s. %v", GithubInfo.DecompressTo, GithubInfo.DecompressedFilesDestination, err))
	//}
	//
	//color.White("- Copy %s to %s success", srcPath, GithubInfo.DecompressedFilesDestination)

	return nil
}

func goGetFromRemoteUrl(ctx *cli.Context) error {
	tag := getTagFromFlags(ctx)
	if len(tag) < 1 {
		tag = "latest"
	}

	common.GetEvent(ctx).AddPayloads(zap.String("tag", tag))
	color.White("- %s@%s", GithubInfo.GoGetUrl, tag)

	_, err := exec.Command("go", "get", "-v", fmt.Sprintf("%s@%s", GithubInfo.GoGetUrl, tag)).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to go get from remote repo. %v", err))
	}

	return nil
}

func validateInstallation(ctx *cli.Context) error {
	bytes, err := GithubInfo.ValidationCmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to validate. %v", err))
	}

	color.White("%s", strings.TrimRight(string(bytes), "\n"))
	return nil
}

func createTempDirForDecompress() error {
	dir, _ := path.Split(GithubInfo.DownloadTo)

	target := path.Join(dir, GithubInfo.Repo+"-"+strconv.Itoa(time.Now().Nanosecond()))
	if err := os.MkdirAll(target, 0777); err != nil {
		return err
	}

	GithubInfo.DecompressTo = target
	return nil
}

// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
	"github.com/schollz/progressbar/v3"
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
)

type URLFilter func(string) bool

type GithubReleaseContext struct {
	Release       *github.RepositoryRelease
	Repo          string
	Owner         string
	Filter        func(string) bool
	RemoteURL     string
	LocalFilePath string
	ExtractPath   string
	ExtractType   string
	ExtractArg    string
	TempPath      string
	DestPath      string
	AssetSize     int64
}

func init() {
	userHomeDir, _ := os.UserHomeDir()
	RkHomeDir = path.Join(userHomeDir, ".rk")
	os.MkdirAll(RkHomeDir, os.ModePerm)
}

func Success() {
	color.Green("[success]")
	color.White("--------------------------------")
}

func Error(event rk_query.Event, err error) error {
	color.Red(err.Error())
	rk_common.Finish(event, err)
	return err
}

func PrintReleasesFromGithub(owner, repo string, event rk_query.Event) error {
	event.AddPair("owner", SwagOwner)
	event.AddPair("repo", SwagRepo)
	event.SetCounter("list_page_size", 5)

	opt := &github.ListOptions{
		PerPage: 5,
	}

	res, _, err := GithubClient.Repositories.ListReleases(context.Background(), owner, repo, opt)

	if err != nil {
		return Error(event, rk_common.NewGithubClientError(
			fmt.Sprintf("failed to list repo:%s release from github\n[err] %v", repo, err)))
	}

	if InstallInfo.Debug {
		printRelease(res...)
	} else {
		for i := range res {
			color.Cyan("%s", res[i].GetTagName())
		}
	}

	if len(res) < 1 {
		color.Cyan("no release found, please use default release")
	}

	return nil
}

func GetReleaseContext(owner, repo, release string, event rk_query.Event) (*GithubReleaseContext, error) {
	event.AddPair("owner", SwagOwner)
	event.AddPair("repo", SwagRepo)

	if len(release) < 1 {
		color.Cyan("Search release name with release:latest")
	}

	opt := &github.ListOptions{PerPage: 1}

	ctx := &GithubReleaseContext{
		Owner: owner,
		Repo:  repo,
		// this should be set it negative one while create a new one
		AssetSize: -1,
		DestPath:  UserLocalBin,
	}

	if len(release) < 1 {
		// list releases
		releases, _, err := GithubClient.Repositories.ListReleases(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, Error(event,
				rk_common.NewGithubClientError(
					fmt.Sprintf("failed to get %s:%s latest release from github\n[err] %v", ctx.Owner, ctx.Repo, err)))
		}

		if len(releases) < 1 {
			return nil, Error(event,
				rk_common.NewGithubClientError(
					fmt.Sprintf("failed to get %s:%s latest release from github\n[err] %v", ctx.Owner, ctx.Repo, err)))
		}

		if InstallInfo.Debug {
			printRelease(releases...)
		}

		// get latest and assign to context
		ctx.Release = releases[0]
	} else {
		var res *github.RepositoryRelease
		var err error
		if res, _, err = GithubClient.Repositories.GetReleaseByTag(context.Background(), owner, repo, release); err != nil {
			return nil, Error(event,
				rk_common.NewGithubClientError(
					fmt.Sprintf("failed to get release:%s from github", release)))
		}

		ctx.Release = res
	}

	color.Cyan("Get context with release:%s", ctx.Release.GetTagName())
	// get release by latest release
	res, _, err := GithubClient.Repositories.GetReleaseByTag(context.Background(), ctx.Owner, ctx.Repo, ctx.Release.GetTagName())
	if err != nil {
		return nil, Error(event,
			rk_common.NewGithubClientError(
				fmt.Sprintf("failed to get assets %s %s from github\n[err] %v", ctx.Repo, ctx.Release.GetTagName(), err)))
	}

	if InstallInfo.Debug {
		printRelease(res)
	}

	if res == nil || len(res.Assets) < 1 {
		return nil, Error(event,
			rk_common.NewGithubClientError(
				fmt.Sprintf("assets is empty %s %s from github\n[err] %v", ctx.Repo, ctx.Release.GetName(), err)))
	}

	return ctx, nil
}

func DownloadGithubRelease(ctx *GithubReleaseContext, event rk_query.Event) error {
	if ctx.Filter == nil {
		return Error(event,
			rk_common.NewInvalidParamError(
				fmt.Sprintf("download filter can not be empty")))
	}

	// get id based on filter
	var id = int64(-1)
	var downloadURL = ""
	for i := range ctx.Release.Assets {
		url := *ctx.Release.Assets[i].BrowserDownloadURL
		if ctx.Filter(url) {
			id = ctx.Release.Assets[i].GetID()
			downloadURL = url
			ctx.AssetSize = int64(ctx.Release.Assets[i].GetSize())
			event.AddPair("asset_id", strconv.FormatInt(id, 10))
			event.AddPair("asset_size_bytes", strconv.FormatInt(ctx.AssetSize, 10))
			break
		}
	}

	ctx.RemoteURL = downloadURL
	ctx.LocalFilePath = path.Join(RkHomeDir, path.Base(downloadURL))
	event.AddPair("local_file_path", ctx.LocalFilePath)
	color.Cyan("Download %s to %s", ctx.RemoteURL, ctx.LocalFilePath)

	if id == -1 {
		return Error(event,
			rk_common.NewInvalidParamError(
				fmt.Sprintf("failed to get downloadable id with release:%s", ctx.Release.GetName())))
	}

	// get reader object
	reader, _, err := GithubClient.Repositories.DownloadReleaseAsset(
		context.Background(),
		ctx.Owner,
		ctx.Repo,
		id,
		http.DefaultClient)

	if err != nil {
		return Error(event,
			rk_common.NewGithubClientError(
				fmt.Sprintf("failed to call github API for download Assert id:%d\n[err] %v", id, err)))
	}

	// create and open/create local file based on download URL
	f, err := os.OpenFile(ctx.LocalFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to open file at %s\n[err] %v", ctx.LocalFilePath, err)))
	}
	defer f.Close()

	// start downloading
	bar := progressbar.DefaultBytes(ctx.AssetSize, "downloading")
	_, err = io.Copy(io.MultiWriter(f, bar), reader)
	if err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to write to file at %s\n[err] %v", ctx.LocalFilePath, err)))
	}

	if ctx.AssetSize < 0 {
		println()
	}

	event.AddPair("remote_url", ctx.RemoteURL)
	event.AddPair("local_file_path", ctx.LocalFilePath)
	event.AddPair("release", ctx.Release.GetName())
	return nil
}

// extract with tar command
func ExtractToDest(ctx *GithubReleaseContext, event rk_query.Event) error {
	// create temp dir
	if err := createTempDir(ctx, event); err != nil {
		return err
	}
	color.White("temporary directory create at %s success", ctx.TempPath)

	var bytes []byte
	var err error
	if ctx.ExtractType == "tar" {
		if len(ctx.ExtractArg) > 0 {
			bytes, err = exec.Command("tar", "-C", ctx.TempPath, "-xvzf", ctx.LocalFilePath, ctx.ExtractArg).CombinedOutput()
		} else {
			bytes, err = exec.Command("tar", "-C", ctx.TempPath, "-xvzf", ctx.LocalFilePath).CombinedOutput()
		}
	} else if ctx.ExtractType == "zip" {
		if len(ctx.ExtractArg) > 0 {
			bytes, err = exec.Command("unzip", "-o", ctx.LocalFilePath, "-d", ctx.TempPath, ctx.ExtractArg).CombinedOutput()
		} else {
			bytes, err = exec.Command("unzip", "-o", ctx.LocalFilePath, "-d", ctx.TempPath).CombinedOutput()
		}
	} else {
		return Error(event, rk_common.NewInvalidParamError(fmt.Sprintf("invalid ExtractType:%s", ctx.ExtractType)))
	}

	if err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to extract %s to %s\n[err] %v\n[stderr] %s", ctx.LocalFilePath, ctx.TempPath, err, string(bytes))))
	}
	color.White("extract %s to %s success", ctx.LocalFilePath, ctx.TempPath)

	// copy to target folder
	// chmod to protoc
	bytes, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("chmod -R +rx %s", path.Join(ctx.TempPath, ctx.ExtractPath))).CombinedOutput()
	bytes, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo cp -r %s %s", path.Join(ctx.TempPath, ctx.ExtractPath), ctx.DestPath)).CombinedOutput()

	if err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to copy %s to %s\n[err] %v\n[stderr] %s", path.Join(ctx.TempPath, ctx.ExtractPath), ctx.DestPath, err, string(bytes))))
	}
	color.White("copy %s to %s success", path.Join(ctx.TempPath, ctx.ExtractPath), ctx.DestPath)

	event.AddPair("temp_path", ctx.TempPath)
	event.AddPair("dest_path", ctx.DestPath)
	event.AddPair("extract_type", ctx.ExtractType)
	event.AddPair("extract_path", ctx.ExtractPath)

	return nil
}

func GoGetFromGithub(repo, url, release string, event rk_query.Event) error {
	color.Cyan("Install %s@%s", repo, release)
	event.AddPair("repo", repo)
	event.AddPair("url", url)
	event.AddPair("releases", release)

	if len(release) > 1 {
		url += fmt.Sprintf("@%s", release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		return Error(event,
			rk_common.NewGithubClientError(
				fmt.Sprintf("failed to install %s\n[err] %v\n[stderr] %s", repo, err, string(bytes))))
	}

	if InstallInfo.Debug {
		color.Yellow("[stdout] %s", string(bytes))
	}

	return nil
}

func ValidateInstallation(command *exec.Cmd, event rk_query.Event) error {
	color.Cyan("Validate installation")
	bytes, err := command.CombinedOutput()
	if err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to show version\n[err] %v\n[stderr] %s", err, string(bytes))))
	}

	color.White("%s", strings.TrimRight(string(bytes), "\n"))
	color.Green("[success]")
	return nil
}

func printRelease(releases ...*github.RepositoryRelease) {
	type releaseSimple struct {
		TagName     *string           `json:"tag_name,omitempty"`
		PublishedAt *github.Timestamp `json:"published_at,omitempty"`
		Assets      []struct {
			Name               *string `json:"name,omitempty"`
			BrowserDownloadURL *string `json:"browser_download_url,omitempty"`
		} `json:"assets,omitempty"`
	}

	for i := range releases {
		release := releases[i]
		var simple releaseSimple

		// marshal to complete json string
		bytes, _ := json.Marshal(release)

		// unmarshal string to simple one
		json.Unmarshal(bytes, &simple)

		// marshal it again with indent
		bytes, _ = json.MarshalIndent(simple, "", "  ")

		// print it
		color.Yellow(fmt.Sprintf(string(bytes)))
	}

}

func createTempDir(ctx *GithubReleaseContext, event rk_query.Event) error {
	// split file name and directory
	dir, _ := path.Split(ctx.LocalFilePath)

	target := path.Join(dir, ctx.Repo+"-"+strconv.Itoa(time.Now().Nanosecond()))
	if err := os.MkdirAll(target, 0777); err != nil {
		return Error(event,
			rk_common.NewFileOperationError(
				fmt.Sprintf("failed to create target folder %s\n[err] %v", target, err)))
	}
	ctx.TempPath = target

	return nil
}

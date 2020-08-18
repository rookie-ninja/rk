package rk_install

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk/common"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	ProtobufOwner = "protocolbuffers"
	ProtobufRepo  = "protobuf"
)

type installProtobufInfo struct {
	ListReleases bool
	Release      string
}

var InstallProtobufInfo = installProtobufInfo{
	ListReleases: false,
	Release:      "latest",
}

// Install protobuf on target hosts
func InstallProtobufCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protobuf",
		Usage:     "install protobuf on local machine",
		UsageText: "rk install protobuf -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallProtobufInfo.Release,
				Usage:       "protobuf release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallProtobufInfo.ListReleases,
				Usage:       "list protobuf releases, list most recent 10 releases",
			},
		},
		Action: InstallProtobufAction,
	}

	return command
}

func InstallProtobufAction(ctx *cli.Context) error {
	if InstallProtobufInfo.ListReleases {
		return listProtobufReleases()
	}

	// get package id
	color.Cyan("Get release id of release:%s", InstallProtobufInfo.Release)
	id, release, err := getProtobufReleaseId(InstallProtobufInfo.Release)
	if err != nil {
		return err
	}
	color.Green("[Success]")

	// download from github
	pkgName := fmt.Sprintf("protoc-%s-%s.zip", release, rk_common.GetSysAndArch())
	color.Cyan("Download %s to %s", pkgName, rk_common.RkHomeDir)
	pkgPath, err := downloadProtobuf(pkgName, id)
	if err != nil {
		return err
	}
	color.Green("[Success]")

	// install
	color.Cyan("Install protobuf at %s", pkgPath)
	err = installProtobuf(pkgPath)
	if err != nil {
		return err
	}
	color.Green("[Success]")

	// show version
	color.Cyan("Validate installation")
	return validateProtobufInstallation()
}

// list available protobuf releases from github by 10 records
func listProtobufReleases() error {
	opt := &github.ListOptions{PerPage: 10}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtobufOwner, ProtobufRepo, opt)
	if err != nil {
		color.Red("failed to list protobuf release from github\n[err] %v", err)
		return err
	}

	for i := range res {
		color.Cyan("%s", res[i].GetTagName())
	}

	if len(res) < 1 {
		color.Cyan("no release found, please use default release")
	}

	return nil
}

// download to local
func downloadProtobuf(fullName string, id int64) (string, error) {
	reader, _, err := rk_common.GithubClient.Repositories.DownloadReleaseAsset(
		context.Background(),
		ProtobufOwner,
		ProtobufRepo,
		id,
		http.DefaultClient)

	if err != nil {
		color.Red("failed to call github API for download Assert id:%d\n[err] %v", id, err)
		return "", err
	}

	pkgPath := path.Join(rk_common.RkHomeDir, fullName)
	f, err := os.OpenFile(pkgPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		color.Red("failed to open file at %s\n[err] %v", pkgPath, err)
		return "", err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(-1, "downloading")
	_, err = io.Copy(io.MultiWriter(f, bar), reader)
	if err != nil {
		if InstallInfo.Debug {
			color.Red("failed to write to file at %s\n[err] %v", pkgPath, err)
		}
	}
	println()
	return pkgPath, nil
}

// install
func installProtobuf(filePath string) error {
	// unzip file
	bytes, err := exec.Command("unzip", "-o", filePath, "-d", "/usr/local", "bin/protoc").CombinedOutput()

	if err != nil {
		color.Red("failed to unzip bin/proto to /usr/local\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	bytes, err = exec.Command("unzip", "-o", filePath, "-d", "/usr/local", "include/*").CombinedOutput()

	if err != nil {
		color.Red("failed to unzip include/* to /usr/local\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	return nil
}

// show version
func validateProtobufInstallation() error {
	bytes, err := exec.Command("protoc", "--version").CombinedOutput()
	if err != nil {
		color.Red("failed to show protobuf version\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

// get latest
func getLatestProtobufRelease() (string, error) {
	opt := &github.ListOptions{PerPage: 1}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtobufOwner, ProtobufRepo, opt)
	if err != nil {
		color.Red("failed to get protobuf latest release from github\n[err] %v", err)
		return "", err
	}

	if len(res) < 1 {
		color.Red("empty release from github")
		return "", errors.New("empty release")
	}

	return res[0].GetName(), nil
}

// get release id with release name
func getProtobufReleaseId(release string) (int64, string, error) {
	if len(release) < 1 || release == "latest" {
		var err error
		release, err = getLatestProtobufRelease()
		if err != nil {
			return -1, "", err
		}
	}

	res, _, err := rk_common.GithubClient.Repositories.GetReleaseByTag(context.Background(), ProtobufOwner, ProtobufRepo, release)

	if err != nil {
		color.Red("failed to list protobuf release from github\n[err] %v", err)
		return -1, "", err
	}

	sysArch := rk_common.GetSysAndArch()
	for i := range res.Assets {
		url := *(res.Assets[i].BrowserDownloadURL)
		if strings.Contains(url, sysArch) {
			return res.Assets[i].GetID(), release, nil
		}
	}

	return -1, "", nil
}

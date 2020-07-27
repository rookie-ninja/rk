package rk_install

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk-cmd/common"
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
	ProtobufRepo = "protobuf"
)

// Install golang on target hosts
func InstallProtobufCommand() *cli.Command{
	command := &cli.Command{
		Name: "protobuf",
		Aliases: []string{"pb"},
		Usage: "protobuf",
		UsageText: "rk install protobuf -l",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "version",
				Aliases: []string{"v"},
				Value: "",
				Usage: "install with version",
			},
		},
		Action: InstallProtobufAction,
	}

	return command
}

func InstallProtobufAction(ctx *cli.Context) error {
	var version = ""

	if version = ctx.String("version"); len(version) < 1 {
		color.Red("protobuf version should be provided with -v xxx\nsee bellow releases:")
		color.Yellow(listProtobufPkg())
		return nil
	}

	sysArch := rk_common.GetSysAndArch()

	// get package id
	color.Cyan("1: Get package id with version:%s, sysArch:%s", version, sysArch)
	id, err := getProtobufPkgId(version, sysArch)
	if err != nil {
		color.Red("err: %v", err)
		return err
	}
	color.Green("Done")

	// download from github
	pkgName := fmt.Sprintf("protoc-%s-%s.zip", version, sysArch)
	color.Cyan("2: Download %s to %s", pkgName, rk_common.RkHomeDir)
	pkgPath, err := downloadProtobufPkg(pkgName, id)
	if err != nil {
		color.Red("err: %v", err)
		return err
	}
	color.Green("Done")

	// install pkg
	color.Cyan("3: Install package %s", pkgPath)
	err = installProtobufPkg(pkgPath)
	if err != nil {
		color.Red("err: %v", err)
		return err
	}
	color.Green("Done")

	return nil
}

func listProtobufPkg() string {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(),
		ProtobufOwner, ProtobufRepo, opt)
	if err != nil {
		return fmt.Sprintf("failed to list protobuf release from github\n %v", err)
	}

	buf := &bytes.Buffer{}
	for i := range res {
		buf.WriteString(res[i].GetTagName() + "\n")
	}

	return buf.String()
}

func downloadProtobufPkg(protobufPkgName string, id int64) (string, error) {
	reader, _, err := rk_common.GithubClient.Repositories.DownloadReleaseAsset(
		context.Background(),
		ProtobufOwner,
		ProtobufRepo,
		id,
		http.DefaultClient)

	if err != nil {
		return "", err
	}

	pkgPath := path.Join(rk_common.RkHomeDir, protobufPkgName)
	f, err := os.OpenFile(pkgPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(-1, "downloading")
	io.Copy(io.MultiWriter(f, bar), reader)

	return pkgPath, nil
}

func installProtobufPkg(filePath string) error {
	// unzip file
	_, err := exec.Command("unzip", "-o", filePath, "-d", "/usr/local", "bin/protoc").CombinedOutput()

	if err != nil {
		return err
	}

	_, err = exec.Command("unzip", "-o", filePath, "-d", "/usr/local", "include/*").CombinedOutput()

	if err != nil {
		return err
	}

	return nil
}

func getProtobufPkgId(version, sysArch string) (int64, error) {
	res, _, err := rk_common.GithubClient.Repositories.GetReleaseByTag(context.Background(), ProtobufOwner, ProtobufRepo, version)

	if err != nil {
		return 0, err
	}

	for i:= range res.Assets {
		url := *(res.Assets[i].BrowserDownloadURL)
		if strings.Contains(url, sysArch) {
			return res.Assets[i].GetID(), nil
		}
	}

	return -1, nil
}



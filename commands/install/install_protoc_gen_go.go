package rk_install

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	rk_common "github.com/rookie-ninja/rk-cmd/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

const (
	ProtocGenGoOwner = "golang"
	ProtocGenGoRepo  = "protobuf"
	ProtocGenGoUrlBase = "github.com/golang/protobuf/protoc-gen-go"
)

type installProtocGenGoInfo struct {
	ListReleases bool
	Release string
}

var InstallProtocGenGoInfo = installProtocGenGoInfo{
	ListReleases: false,
	Release: "",
}

// Install protoc-gen-go on target hosts
func InstallProtocGenGoCommand() *cli.Command{
	command := &cli.Command{
		Name: "protoc-gen-go",
		Usage: "install protoc-gen-go on local machine",
		UsageText: "rk install protoc-gen-go -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallProtocGenGoInfo.Release,
				Required: false,
				Usage:       "protoc-gen-go release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallProtocGenGoInfo.ListReleases,
				Usage:       "list protoc-gen-go releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenGoAction,
	}

	return command
}

func listProtocGenGoReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtocGenGoOwner, ProtocGenGoRepo, opt)
	if err != nil {
		color.Red("failed to list protoc-gen-go release from github\n[err] %v", err)
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

func InstallProtocGenGoAction(ctx *cli.Context) error {
	if InstallProtocGenGoInfo.ListReleases {
		return listProtocGenGoReleases()
	}

	// install pkg
	color.Cyan("Install protoc-gen-go %s", InstallProtocGenGoInfo.Release)
	err := installProtocGenGo()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateProtocGenGoInstallation()
}

func installProtocGenGo() error {
	url := fmt.Sprintf(ProtocGenGoUrlBase)

	if len(InstallProtocGenGoInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallProtocGenGoInfo.Release)
	}

	bytes, err := exec.Command("go", "get", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install protoc-gen-go\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	return nil
}

func validateProtocGenGoInstallation() error {
	bytes, err := exec.Command("which", "protoc-gen-go").CombinedOutput()
	if err != nil {
		color.Red("failed to validate protoc-gen-go\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}
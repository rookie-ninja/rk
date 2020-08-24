package rk_install

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

const (
	ProtocGenDocOwner   = "pseudomuto"
	ProtocGenDocRepo    = "protoc-gen-doc"
	ProtocGenDocUrlBase = "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
)

type installProtocGenDocInfo struct {
	ListReleases bool
	Release      string
}

var InstallProtocGenDocInfo = installProtocGenDocInfo{
	ListReleases: false,
	Release:      "",
}

// Install protoc-gen-grpc-gateway on target hosts
func InstallProtocGenDocCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-doc",
		Usage:     "install protoc-gen-doc on local machine",
		UsageText: "rk install protoc-gen-doc -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallProtocGenDocInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-doc release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallProtocGenDocInfo.ListReleases,
				Usage:       "list protoc-gen-doc releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenDocAction,
	}

	return command
}

func listProtocGenDocReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtocGenDocOwner, ProtocGenDocRepo, opt)
	if err != nil {
		color.Red("failed to list protoc-gen-doc release from github\n[err] %v", err)
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

func InstallProtocGenDocAction(ctx *cli.Context) error {
	if InstallProtocGenDocInfo.ListReleases {
		return listProtocGenDocReleases()
	}

	// install pkg
	color.Cyan("Install protoc-gen-doc %s", InstallProtocGenDocInfo.Release)
	err := installProtocGenDoc()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateProtocGenDocInstallation()
}

func installProtocGenDoc() error {
	url := fmt.Sprintf(ProtocGenDocUrlBase)

	if len(InstallProtocGenDocInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallProtocGenDocInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install protoc-gen-doc\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateProtocGenDocInstallation() error {
	bytes, err := exec.Command("protoc-gen-doc", "-version").CombinedOutput()
	if err != nil {
		color.Red("failed to validate protoc-gen-doc\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

package rk_install

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk-cmd/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

const (
	GoLangCILintOwner = "golangci"
	GoLangCILintRepo  = "golangci-lint"
	GoLangCILintUrlBase = "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

type installGoLangCILintInfo struct {
	ListReleases bool
	Release string
}

var InstallGoLangCILintInfo = installGoLangCILintInfo{
	ListReleases: false,
	Release: "",
}

// Install golangci-lint on target hosts
func InstallGoLangCILintCommand() *cli.Command{
	command := &cli.Command{
		Name: "golangci-lint",
		Usage: "install golangci-lint on local machine",
		UsageText: "rk install golangci-lint -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallGoLangCILintInfo.Release,
				Required: false,
				Usage:       "golangci-lint release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallGoLangCILintInfo.ListReleases,
				Usage:       "list golangci-lint releases, list most recent 10 releases",
			},
		},
		Action: InstallGoLangCILintAction,
	}

	return command
}

func listGoLangCILintReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), GoLangCILintOwner, GoLangCILintRepo, opt)
	if err != nil {
		color.Red("failed to list golangci-lint release from github\n[err] %v", err)
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

func InstallGoLangCILintAction(ctx *cli.Context) error {
	if InstallGoLangCILintInfo.ListReleases {
		return listGoLangCILintReleases()
	}

	// install pkg
	color.Cyan("Install golangci-lint %s", InstallGoLangCILintInfo.Release)
	err := installGoLangCILint()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateGoLangCILintInstallation()
}

func installGoLangCILint() error {
	url := fmt.Sprintf(GoLangCILintUrlBase)

	if len(InstallGoLangCILintInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallGoLangCILintInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install golangci-lint\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateGoLangCILintInstallation() error {
	bytes, err := exec.Command("golangci-lint", "--version").CombinedOutput()
	if err != nil {
		color.Red("failed to validate golangci-lint\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

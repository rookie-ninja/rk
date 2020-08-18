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
	GoLintOwner = "golang"
	GoLintRepo  = "lint"
	GoLintUrlBase = "golang.org/x/lint/golint"
)

type installGoLintInfo struct {
	ListReleases bool
	Release string
}

var InstallGoLintInfo = installGoLintInfo{
	ListReleases: false,
	Release: "",
}

// Install protoc-gen-swagger on target hosts
func InstallGoLintCommand() *cli.Command{
	command := &cli.Command{
		Name: "golint",
		Usage: "install golint on local machine",
		UsageText: "rk install golint -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallGoLintInfo.Release,
				Required: false,
				Usage:       "golint release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallGoLintInfo.ListReleases,
				Usage:       "list golint releases, list most recent 10 releases",
			},
		},
		Action: InstallGoLintAction,
	}

	return command
}

func listGoLintReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), GoLintOwner, GoLintRepo, opt)
	if err != nil {
		color.Red("failed to list golint release from github\n[err] %v", err)
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

func InstallGoLintAction(ctx *cli.Context) error {
	if InstallGoLintInfo.ListReleases {
		return listGoLintReleases()
	}

	// install pkg
	color.Cyan("Install golint %s", InstallGoLintInfo.Release)
	err := installGoLint()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateGoLintInstallation()
}

func installGoLint() error {
	url := fmt.Sprintf(GoLintUrlBase)

	if len(InstallGoLintInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallGoLintInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install golint\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateGoLintInstallation() error {
	bytes, err := exec.Command("which", "golint").CombinedOutput()
	if err != nil {
		color.Red("failed to validate golint\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

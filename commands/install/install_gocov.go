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
	GoCovOwner   = "axw"
	GoCovRepo    = "gocov"
	GoCovUrlBase = "github.com/axw/gocov/gocov"
)

type installGoCovInfo struct {
	ListReleases bool
	Release      string
}

var InstallGoCovInfo = installGoCovInfo{
	ListReleases: false,
	Release:      "",
}

// Install gocov on target hosts
func InstallGoCovCommand() *cli.Command {
	command := &cli.Command{
		Name:      "gocov",
		Usage:     "install gocov on local machine",
		UsageText: "rk install gocov -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallGoCovInfo.Release,
				Required:    false,
				Usage:       "gocov release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallGoCovInfo.ListReleases,
				Usage:       "list gocov releases, list most recent 10 releases",
			},
		},
		Action: InstallGoCovAction,
	}

	return command
}

func listGoCovReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), GoCovOwner, GoCovRepo, opt)
	if err != nil {
		color.Red("failed to list gocov release from github\n[err] %v", err)
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

func InstallGoCovAction(ctx *cli.Context) error {
	if InstallGoCovInfo.ListReleases {
		return listGoCovReleases()
	}

	// install pkg
	color.Cyan("Install gocov %s", InstallGoCovInfo.Release)
	err := installGoCov()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateGoCovInstallation()
}

func installGoCov() error {
	url := fmt.Sprintf(GoCovUrlBase)

	if len(InstallGoCovInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallGoCovInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install gocov\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateGoCovInstallation() error {
	bytes, err := exec.Command("which", "gocov").CombinedOutput()
	if err != nil {
		color.Red("failed to validate gocov\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

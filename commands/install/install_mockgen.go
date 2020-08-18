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
	MockGenOwner = "golang"
	MockGenRepo  = "mock"
	MockGenUrlBase = "github.com/golang/mock/mockgen"
)

type installMockGenInfo struct {
	ListReleases bool
	Release string
}

var InstallMockGenInfo = installMockGenInfo{
	ListReleases: false,
	Release: "",
}

// Install gocov on target hosts
func InstallMockGenCommand() *cli.Command{
	command := &cli.Command{
		Name: "mockgen",
		Usage: "install mockgen on local machine",
		UsageText: "rk install mockgen -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallMockGenInfo.Release,
				Required: false,
				Usage:       "mockgen release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallMockGenInfo.ListReleases,
				Usage:       "list mockgen releases, list most recent 10 releases",
			},
		},
		Action: InstallMockGenAction,
	}

	return command
}

func listMockGenReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), MockGenOwner, MockGenRepo, opt)
	if err != nil {
		color.Red("failed to list mockgen release from github\n[err] %v", err)
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

func InstallMockGenAction(ctx *cli.Context) error {
	if InstallMockGenInfo.ListReleases {
		return listMockGenReleases()
	}

	// install pkg
	color.Cyan("Install mockgen %s", InstallMockGenInfo.Release)
	err := installMockGen()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateMockGenInstallation()
}

func installMockGen() error {
	url := fmt.Sprintf(MockGenUrlBase)

	if len(InstallMockGenInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallMockGenInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install mockgen\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateMockGenInstallation() error {
	bytes, err := exec.Command("which", "mockgen").CombinedOutput()
	if err != nil {
		color.Red("failed to validate mockgen\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

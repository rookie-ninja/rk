package rk_install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

const (
	MockGenOwner   = "golang"
	MockGenRepo    = "mock"
	MockGenUrlBase = "github.com/golang/mock/mockgen"
)

// Install mockgen on target hosts with go get command
func InstallMockGenCommand() *cli.Command {
	command := &cli.Command{
		Name:      "mockgen",
		Usage:     "install mockgen on local machine",
		UsageText: "rk install mockgen -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "mockgen release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list mockgen releases, list most recent 10 releases",
			},
		},
		Action: InstallMockGenAction,
	}

	return command
}

func InstallMockGenAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-mock-gen-release")
		return PrintReleasesFromGithub(MockGenOwner, MockGenRepo, event)
	}

	event := rk_common.GetEvent("install-mock-gen")
	// go get package
	if err := GoGetFromGithub(MockGenRepo, MockGenUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("which", "mockgen"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

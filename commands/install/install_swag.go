// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"runtime"
	"strings"
)

const (
	SwagOwner = "swaggo"
	SwagRepo  = "swag"
)

// Install swag on target hosts
func InstallSwagCommand() *cli.Command {
	command := &cli.Command{
		Name:      "swag",
		Usage:     "install swag on local machine",
		UsageText: "rk install swag -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "swag release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list swag releases, list most recent 10 releases",
			},
		},
		Action: InstallSwagAction,
	}

	return command
}

func InstallSwagAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-swag-release")
		return PrintReleasesFromGithub(SwagOwner, SwagRepo, event)
	}

	event := rk_common.GetEvent("install-swag")
	// 1: get release context
	releaseCtx, err := GetReleaseContext(SwagOwner, SwagRepo, InstallInfo.Release, event)
	if err != nil {
		return err
	}
	Success()

	// 2.1: specify filter based on pattern
	// https://github.com/swaggo/swag/releases
	releaseCtx.Filter = func(url string) bool {
		var sysArch string
		if runtime.GOOS == "darwin" {
			sysArch = "Darwin_x86_64"
		} else {
			sysArch = "Linux_x86_64"
		}

		if strings.Contains(url, sysArch) {
			return true
		}

		return false
	}

	// 2.2: download with filter
	DownloadGithubRelease(releaseCtx, event)
	if err != nil {
		return err
	}
	Success()

	// 3: install release
	color.Cyan("Install swag in %s to %s", releaseCtx.LocalFilePath, UserLocalBin)
	releaseCtx.ExtractType = "tar"
	releaseCtx.ExtractPath = "swag"
	if err = ExtractToDest(releaseCtx, event); err != nil {
		return err
	}
	Success()

	// 4: validate installation
	if err := ValidateInstallation(exec.Command("swag", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

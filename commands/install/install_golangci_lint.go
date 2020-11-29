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
	GoLangCILintOwner = "golangci"
	GoLangCILintRepo  = "golangci-lint"
)

// Install golangci-lint on target hosts directly download assets from github
func InstallGoLangCILintCommand() *cli.Command {
	command := &cli.Command{
		Name:      "golangci-lint",
		Usage:     "install golangci-lint on local machine",
		UsageText: "rk install golangci-lint -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "golangci-lint release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list golangci-lint releases, list most recent 10 releases",
			},
		},
		Action: InstallGoLangCILintAction,
	}

	return command
}

func InstallGoLangCILintAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-golangci-lint-release")
		return PrintReleasesFromGithub(GoLangCILintOwner, GoLangCILintRepo, event)
	}

	event := rk_common.GetEvent("install-golangci-lint")
	// 1: get release context
	releaseCtx, err := GetReleaseContext(GoLangCILintOwner, GoLangCILintRepo, InstallInfo.Release, event)
	if err != nil {
		return err
	}
	Success()

	// 2.1: specify filter based on pattern
	// https://github.com/golangci/golangci-lint/releases
	releaseCtx.Filter = func(url string) bool {
		sysArch := runtime.GOOS + "-" + runtime.GOARCH

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
	color.Cyan("install golangci-lint in %s to %s", releaseCtx.LocalFilePath, UserLocalBin)
	releaseCtx.ExtractType = "tar"
	releaseCtx.ExtractArg = "--strip-components=1"
	releaseCtx.ExtractPath = "golangci-lint"
	if err = ExtractToDest(releaseCtx, event); err != nil {
		return err
	}
	Success()

	// 4: validate installation
	if err := ValidateInstallation(exec.Command("golangci-lint", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

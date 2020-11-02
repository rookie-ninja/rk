// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

const (
	GoLintOwner   = "golang"
	GoLintRepo    = "lint"
	GoLintUrlBase = "golang.org/x/lint/golint"
)

// Install golint on target hosts with go get command
func InstallGoLintCommand() *cli.Command {
	command := &cli.Command{
		Name:      "golint",
		Usage:     "install golint on local machine",
		UsageText: "rk install golint -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "golint release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list golint releases, list most recent 10 releases",
			},
		},
		Action: InstallGoLintAction,
	}

	return command
}

func InstallGoLintAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-golint-release")
		return PrintReleasesFromGithub(GoLintOwner, GoLintRepo, event)
	}

	event := rk_common.GetEvent("install-golint")
	// go get package
	if err := GoGetFromGithub(GoLintRepo, GoLintUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}

	if err := ValidateInstallation(exec.Command("which", "golint"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

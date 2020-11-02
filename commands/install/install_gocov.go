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
	GoCovOwner   = "axw"
	GoCovRepo    = "gocov"
	GoCovUrlBase = "github.com/axw/gocov/gocov"
)

// Install gocov on target hosts with go get command
func InstallGoCovCommand() *cli.Command {
	command := &cli.Command{
		Name:      "gocov",
		Usage:     "install gocov on local machine",
		UsageText: "rk install gocov -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "gocov release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list gocov releases, list most recent 10 releases",
			},
		},
		Action: InstallGoCovAction,
	}

	return command
}

func InstallGoCovAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-gocov-release")
		return PrintReleasesFromGithub(GoCovOwner, GoCovRepo, event)
	}

	event := rk_common.GetEvent("install-gocov")
	// go get package
	if err := GoGetFromGithub(GoCovRepo, GoCovUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("which", "gocov"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

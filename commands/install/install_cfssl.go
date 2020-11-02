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
	CfsslOwner   = "cloudflare"
	CfsslRepo    = "cfssl"
	CfsslUrlBase = "github.com/cloudflare/cfssl/cmd/cfssl"
)

// Install cfssl on target hosts
func InstallCfsslCommand() *cli.Command {
	command := &cli.Command{
		Name:      "cfssl",
		Usage:     "install cfssl on local machine",
		UsageText: "rk install cfssl -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "cfssl release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list cfssl releases, list most recent 10 releases",
			},
		},
		Action: InstallCfsslAction,
	}

	return command
}

func InstallCfsslAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-cfssl-release")
		return PrintReleasesFromGithub(CfsslOwner, CfsslRepo, event)
	}

	event := rk_common.GetEvent("install-cfssl")
	// go get package
	if err := GoGetFromGithub(CfsslRepo, CfsslUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("cfssl", "version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

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
	CfsslJsonOwner   = "cloudflare"
	CfsslJsonRepo    = "cfssl"
	CfsslJsonUrlBase = "github.com/cloudflare/cfssl/cmd/cfssljson"
)

// Install cfssljson on target hosts
func InstallCfsslJsonCommand() *cli.Command {
	command := &cli.Command{
		Name:      "cfssljson",
		Usage:     "install cfssljson on local machine",
		UsageText: "rk install cfssljson -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "cfssljson release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list cfssljson releases, list most recent 10 releases",
			},
		},
		Action: InstallCfsslJsonAction,
	}

	return command
}

func InstallCfsslJsonAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-cfssljson-release")
		return PrintReleasesFromGithub(CfsslJsonOwner, CfsslJsonRepo, event)
	}

	event := rk_common.GetEvent("install-cfssljson")
	// go get package
	if err := GoGetFromGithub(CfsslJsonRepo, CfsslJsonUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("cfssljson", "-version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

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
	ProtocGenSwaggerOwner   = "grpc-ecosystem"
	ProtocGenSwaggerRepo    = "grpc-gateway"
	ProtocGenSwaggerUrlBase = "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
)

// Install protoc-gen-swagger on target hosts
func InstallProtocGenSwaggerCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-openapiv2",
		Usage:     "install protoc-gen-openapiv2 on local machine",
		UsageText: "rk install protoc-gen-openapiv2 -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-openapiv2 release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protoc-gen-openapiv2 releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenSwaggerAction,
	}

	return command
}

func InstallProtocGenSwaggerAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protoc-gen-openapiv2-release")
		return PrintReleasesFromGithub(ProtocGenSwaggerOwner, ProtocGenSwaggerRepo, event)
	}

	event := rk_common.GetEvent("install-protoc-gen-openapiv2")
	// go get package
	if err := GoGetFromGithub(ProtocGenSwaggerRepo, ProtocGenSwaggerUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("protoc-gen-openapiv2", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

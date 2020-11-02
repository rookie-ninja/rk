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
	ProtocGenSwaggerUrlBase = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
)

// Install protoc-gen-swagger on target hosts
func InstallProtocGenSwaggerCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-swagger",
		Usage:     "install protoc-gen-swagger on local machine",
		UsageText: "rk install protoc-gen-swagger -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-swagger release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protoc-gen-swagger releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenSwaggerAction,
	}

	return command
}

func InstallProtocGenSwaggerAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protoc-gen-swagger-release")
		return PrintReleasesFromGithub(ProtocGenSwaggerOwner, ProtocGenSwaggerRepo, event)
	}

	event := rk_common.GetEvent("install-protoc-gen-swagger")
	// go get package
	if err := GoGetFromGithub(ProtocGenSwaggerRepo, ProtocGenSwaggerUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("protoc-gen-swagger", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

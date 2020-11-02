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
	ProtocGenGrpcGatewayOwner   = "grpc-ecosystem"
	ProtocGenGrpcGatewayRepo    = "grpc-gateway"
	ProtocGenGrpcGatewayUrlBase = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
)

// Install protoc-gen-grpc-gateway on target hosts
func InstallProtocGenGrpcGatewayCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-grpc-gateway",
		Usage:     "install protoc-gen-grpc-gateway on local machine",
		UsageText: "rk install protoc-gen-grpc-gateway -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-grpc-gateway release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protoc-gen-grpc-gateway releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenGrpcGatewayAction,
	}

	return command
}

func InstallProtocGenGrpcGatewayAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protoc-gen-grpc-gateway-release")
		return PrintReleasesFromGithub(ProtocGenGrpcGatewayOwner, ProtocGenGrpcGatewayRepo, event)
	}

	event := rk_common.GetEvent("install-protoc-gen-grpc-gateway")
	// go get package
	if err := GoGetFromGithub(ProtocGenGrpcGatewayRepo, ProtocGenGrpcGatewayUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("protoc-gen-grpc-gateway", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

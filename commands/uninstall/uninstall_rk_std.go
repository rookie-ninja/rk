// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package uninstall

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Uninstall rk-std on target hosts
func UninstallRkStdCommand() *cli.Command {
	command := &cli.Command{
		Name:      "rk-std",
		Usage:     "uninstall rk standard environment on local machine",
		UsageText: "rk uninstall rk-std",
		Action:    UninstallRkStdAction,
	}

	return command
}

func UninstallRkStdAction(ctx *cli.Context) error {
	color.Magenta("Uninstall rk style standard commands...")

	// install protobuf
	color.Magenta("Uninstall protobuf")
	if err := protobufAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-go
	color.Magenta("Uninstall protoc-gen-go")
	if err := protocGenGoAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-grpc-gateway
	color.Magenta("Uninstall protoc-gen-grpc-gateway")
	if err := protocGenGrpcGatewayAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-doc
	color.Magenta("Uninstall protoc-gen-doc")
	if err := protocGenDocAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-openapiv2
	color.Magenta("Uninstall protoc-gen-openapiv2")
	if err := protocGenOpenApiV2Action(ctx); err != nil {
		return err
	}

	// install golangci-lint
	color.Magenta("Uninstall golangci-lint")
	if err := golangCiLintAction(ctx); err != nil {
		return err
	}

	// install gocov
	color.Magenta("Uninstall gocov")
	if err := goCovAction(ctx); err != nil {
		return err
	}

	// install swag
	color.Magenta("Uninstall swag")
	if err := swagAction(ctx); err != nil {
		return err
	}

	// install buf
	color.Magenta("Uninstall buf")
	if err := bufAction(ctx); err != nil {
		return err
	}

	color.Magenta("[OK]")
	return nil
}

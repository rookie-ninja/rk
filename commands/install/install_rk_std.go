// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package install

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Install rk-std on target hosts
func rkStdCommand() *cli.Command {
	command := &cli.Command{
		Name:      "rk-std",
		Usage:     "install rk standard environment on local machine",
		UsageText: "rk install rk-std",
		Action:    rkStdAction,
	}

	return command
}

func rkStdAction(ctx *cli.Context) error {
	color.Magenta("Install rk style standard commands...")

	// install protobuf
	color.Magenta("Install protobuf")
	if err := protobufAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-go
	color.Magenta("Install protoc-gen-go")
	if err := protocGenGoAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-grpc-gateway
	color.Magenta("Install protoc-gen-grpc-gateway")
	if err := protocGenGrpcGatewayAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-doc
	color.Magenta("Install protoc-gen-doc")
	if err := protocGenDocAction(ctx); err != nil {
		return err
	}

	// install protoc-gen-openapiv2
	color.Magenta("Install protoc-gen-openapiv2")
	if err := protocGenOpenApiV2Action(ctx); err != nil {
		return err
	}

	// install golangci-lint
	color.Magenta("Install golangci-lint")
	if err := golangCiLintAction(ctx); err != nil {
		return err
	}

	// install gocov
	color.Magenta("Install gocov")
	if err := goCovAction(ctx); err != nil {
		return err
	}

	// install swag
	color.Magenta("Install swag")
	if err := swagAction(ctx); err != nil {
		return err
	}

	// install buf
	color.Magenta("Install buf")
	if err := bufAction(ctx); err != nil {
		return err
	}

	color.Magenta("[OK]")
	return nil
}

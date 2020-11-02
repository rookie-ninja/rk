// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import (
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const EOE = "------------------------------------------------------------------"

// Install rk-std on target hosts
func InstallRkStdCommand() *cli.Command {
	command := &cli.Command{
		Name:      "rk-std",
		Usage:     "install rk standard environment on local machine",
		UsageText: "rk install rk-std",
		Action:    InstallRkStdAction,
	}

	return command
}

func InstallRkStdAction(ctx *cli.Context) error {
	// install protobuf
	color.Magenta("[1] Install protobuf")
	if err := InstallProtobufAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-go
	color.Magenta("[2] Install protoc-gen-go")
	if err := InstallProtocGenGoAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-grpc-gateway
	color.Magenta("[3] Install protoc-gen-grpc-gateway")
	if err := InstallProtocGenGrpcGatewayAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-doc
	color.Magenta("[4] Install protoc-gen-doc")
	if err := InstallProtocGenDocAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-swagger
	color.Magenta("[5] Install protoc-gen-swagger")
	if err := InstallProtocGenSwaggerAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install mockgen
	color.Magenta("[6] Install mockgen")
	if err := InstallMockGenAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install golangci-lint
	color.Magenta("[7] Install golangci-lint")
	if err := InstallGoLangCILintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install gocov
	color.Magenta("[8] Install gocov")
	if err := InstallGoCovAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	color.Magenta("[9] Install golint")
	if err := InstallGoLintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install swag
	color.Magenta("[10] Install swag")
	if err := InstallSwagAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	return nil
}

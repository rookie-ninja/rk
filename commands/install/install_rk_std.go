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
	color.Magenta("[Action] install protobuf")
	if err := InstallProtobufAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-go
	color.Magenta("[Action] install protoc-gen-go")
	if err := InstallProtocGenGoAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-grpc-gateway
	color.Magenta("[Action] install protoc-gen-grpc-gateway")
	if err := InstallProtocGenGrpcGatewayAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-doc
	color.Magenta("[Action] install protoc-gen-doc")
	if err := InstallProtocGenDocAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install protoc-gen-swagger
	color.Magenta("[Action] install protoc-gen-swagger")
	if err := InstallProtocGenSwaggerAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install mockgen
	color.Magenta("[Action] install mockgen")
	if err := InstallMockGenAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install golangci-lint
	color.Magenta("[Action] install golangci-lint")
	if err := InstallGoLangCILintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install gocov
	color.Magenta("[Action] install gocov")
	if err := InstallGoCovAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install golint
	color.Magenta("[Action] install golint")
	if err := InstallGoLintAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	// install swag
	color.Magenta("[Action] install swag")
	if err := InstallSwagAction(ctx); err != nil {
		return err
	}
	color.Magenta("[success]")
	color.White(EOE)

	return nil
}

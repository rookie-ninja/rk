// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package pkg

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

const (
	goGrpcDescShort = `- [go-grpc]: golang microservice package with [gRPC] as framework and with rk-boot as bootstrapper.`
	goGinSDescShort = `- [go-gin]: golang microservice package with [Gin] as framework and with rk-boot as bootstrapper.`
)

// PkgList List available package can be created.
func PkgList() *cli.Command {
	command := &cli.Command{
		Name:      "list",
		Usage:     "List available package types can be created",
		UsageText: "rk pkg list",
		Action:    pkgListAction,
	}

	return command
}

// Action of list sub command
func pkgListAction(ctx *cli.Context) error {
	chain := common.NewActionChain()
	chain.Add("List available package types to crate", func(context *cli.Context) error {
		color.White(goGrpcDescShort)
		color.White(goGinSDescShort)
		return nil
	}, true)

	return chain.Execute(ctx)
}

// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

const (
	goGrpcDescFull = `- [go-grpc]
  Lang: Golang v1.16
  GoMod: Yes
  Framework: gRPC(https://github.com/grpc/grpc-go)
  Dependency: rk-boot, grpc-gateway, grpc, protobuf
  Tools required (Please use rk cli to install bellow tools):
   - doctoc
   - golangci-lint
   - protoc, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-grpc, protoc-gen-grpc-gateway, protoc-gen-openapiv2
   - buf
  Source: https://github.com/rookie-ninja/rk-demo/tree/master/standard/go-grpc
  License: Apache 2.0
  Docker supported: Yes
`
	goGinSDescFull = `- [go-gin]
  Lang: Golang v1.16
  GoMod: Yes
  Framework: Gin(https://github.com/gin-gonic/gin)
  Dependency: rk-boot, grpc-gateway, grpc, protobuf
  Tools required (Please use rk cli to install bellow tools):
   - doctoc
   - golangci-lint
   - swag
  Source: https://github.com/rookie-ninja/rk-demo/tree/master/standard/go-gin
  License: Apache 2.0
  Docker supported: Yes
`
)

var descMap = map[string]string{
	"go-gin":  goGinSDescFull,
	"go-grpc": goGrpcDescFull,
}

// PkgList List available package can be created.
func PkgDesc() *cli.Command {
	command := &cli.Command{
		Name:      "desc",
		Usage:     "Describe packages can be created",
		UsageText: "rk pkg desc -t [type]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "specify which package type hope to describe, all packages will be listed if not provided",
			},
		},
		Action: pkgDescAction,
	}

	return command
}

// Action of list sub command
func pkgDescAction(ctx *cli.Context) error {
	pkgType := getPkgTypeFromFlags(ctx)
	chain := common.NewActionChain()

	pkgTypeForTitle := pkgType
	if len(pkgType) < 1 {
		pkgTypeForTitle = "for all"
	}
	chain.Add(fmt.Sprintf("Describe package %s", pkgTypeForTitle), func(context *cli.Context) error {
		if len(pkgType) < 1 {
			for _, v := range descMap {
				color.White(v)
			}
		} else if v, ok := descMap[pkgType]; ok {
			color.White(v)
		} else {
			color.Yellow("Unknown package type [%s]", pkgType)
		}
		return nil
	}, true)

	return chain.Execute(ctx)
}

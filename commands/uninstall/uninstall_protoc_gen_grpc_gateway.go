// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install protoc-gen-grpc-gateway on target hosts
func UninstallProtocGenGrpcGatewayCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-grpc-gateway",
		Usage:     "uninstall protoc-gen-grpc-gateway on local machine",
		UsageText: "rk uninstall protoc-gen-grpc-gateway",
		Action:    UninstallMockGenAction,
	}

	return command
}

func UninstallProtocGenGrpcGatewayAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-protoc-gen-grpc-gateway")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("protoc-gen-grpc-gateway", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall protoc-gen-grpc-gateway at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

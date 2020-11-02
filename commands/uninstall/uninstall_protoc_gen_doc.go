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

// Install protoc-gen-doc on target hosts
func UninstallProtocGenDocCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-doc",
		Usage:     "uninstall protoc-gen-doc on local machine",
		UsageText: "rk uninstall protoc-gen-doc",
		Action:    UninstallProtocGenDocAction,
	}

	return command
}

func UninstallProtocGenDocAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-protoc-gen-doc")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("protoc-gen-doc", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall protoc-gen-doc at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

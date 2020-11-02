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

// Install swag on target hosts
func UninstallSwagCommand() *cli.Command {
	command := &cli.Command{
		Name:      "swag",
		Usage:     "uninstall swag on local machine",
		UsageText: "rk uninstall swag",
		Action:    UninstallSwagAction,
	}

	return command
}

func UninstallSwagAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-swag")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("swag", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall swag at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

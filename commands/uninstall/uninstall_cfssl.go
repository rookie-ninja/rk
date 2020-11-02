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

// Install cfssl on target hosts
func UninstallCfsslCommand() *cli.Command {
	command := &cli.Command{
		Name:      "cfssl",
		Usage:     "uninstall cfssl on local machine",
		UsageText: "rk uninstall cfssl",
		Action:    UninstallCfsslAction,
	}

	return command
}

func UninstallCfsslAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-cfssl")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("cfssl", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall cfssl at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

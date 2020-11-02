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

// Install cfssljson on target hosts
func UninstallCfsslJsonCommand() *cli.Command {
	command := &cli.Command{
		Name:      "cfssljson",
		Usage:     "uninstall cfssljson on local machine",
		UsageText: "rk uninstall cfssljson",
		Action:    UninstallCfsslJsonAction,
	}

	return command
}

func UninstallCfsslJsonAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-cfssljson")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("cfssljson", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall cfssljson at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

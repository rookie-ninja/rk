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

// Install mockgen on target hosts
func UninstallMockGenCommand() *cli.Command {
	command := &cli.Command{
		Name:      "golangci-lint",
		Usage:     "uninstall golangci-lint on local machine",
		UsageText: "rk uninstall golangci-lint",
		Action:    UninstallMockGenAction,
	}

	return command
}

func UninstallMockGenAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-mockgen")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("mockgen", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall mockgen at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

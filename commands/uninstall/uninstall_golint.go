// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_uninstall

import (
	"fmt"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

func uninstallGolint() *cli.Command {
	command := commandDefault("golint")
	command.Before = beforeDefault
	command.Action = golintAction
	command.After = afterDefault

	return command
}

func golintAction(ctx *cli.Context) error {
	UninstallInfo.app = "golint"

	chain := rk_common.NewActionChain()
	chain.Add(fmt.Sprintf("Check path of %s", UninstallInfo.app), checkPath, false)
	chain.Add("Validate uninstallation", validateUninstallation, false)
	err := chain.Execute(ctx)

	return err
}

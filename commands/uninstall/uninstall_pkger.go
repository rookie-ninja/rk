// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package uninstall

import (
	"fmt"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

func uninstallPkger() *cli.Command {
	command := commandDefault("pkger")
	command.Before = beforeDefault
	command.Action = pkgerAction
	command.After = afterDefault

	return command
}

func pkgerAction(ctx *cli.Context) error {
	UninstallInfo.app = "pkger"

	chain := common.NewActionChain()
	chain.Add(fmt.Sprintf("Check path of %s", UninstallInfo.app), checkPath, false)
	chain.Add("Validate uninstallation", validateUninstallation, false)
	err := chain.Execute(ctx)

	return err
}

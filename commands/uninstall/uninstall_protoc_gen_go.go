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

func uninstallProtocGenGo() *cli.Command {
	command := commandDefault("protoc-gen-go")
	command.Before = beforeDefault
	command.Action = protocGenGoAction
	command.After = afterDefault

	return command
}

func protocGenGoAction(ctx *cli.Context) error {
	UninstallInfo.app = "protoc-gen-go"

	chain := common.NewActionChain()
	chain.Add(fmt.Sprintf("Check path of %s", UninstallInfo.app), checkPath, false)
	chain.Add("Validate uninstallation", validateUninstallation, false)
	err := chain.Execute(ctx)

	return err
}

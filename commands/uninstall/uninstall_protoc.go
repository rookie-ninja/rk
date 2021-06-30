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

func uninstallProtobuf() *cli.Command {
	command := commandDefault("protobuf")
	command.Before = beforeDefault
	command.Action = protobufAction
	command.After = afterDefault

	return command
}

func protobufAction(ctx *cli.Context) error {
	UninstallInfo.app = "protobuf"

	chain := common.NewActionChain()
	chain.Add(fmt.Sprintf("Check path of %s", UninstallInfo.app), checkPath, false)
	chain.Add("Validate uninstallation", validateUninstallation, false)
	err := chain.Execute(ctx)

	return err
}

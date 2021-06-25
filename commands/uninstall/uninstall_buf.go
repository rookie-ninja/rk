package rk_uninstall

import (
	"fmt"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

func uninstallBuf() *cli.Command {
	command := commandDefault("buf")
	command.Before = beforeDefault
	command.Action = bufAction
	command.After = afterDefault

	return command
}

func bufAction(ctx *cli.Context) error {
	UninstallInfo.app = "buf"

	chain := rk_common.NewActionChain()
	chain.Add(fmt.Sprintf("Check path of %s", UninstallInfo.app), checkPath, false)
	chain.Add("Validate uninstallation", validateUninstallation, false)
	err := chain.Execute(ctx)

	return err
}

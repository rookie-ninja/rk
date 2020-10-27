package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install golint on target hosts
func UninstallGoLintCommand() *cli.Command {
	command := &cli.Command{
		Name:      "golint",
		Usage:     "uninstall golint on local machine",
		UsageText: "rk uninstall golint",
		Action:    UninstallGoLintAction,
	}

	return command
}

func UninstallGoLintAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-golint")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("golint", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall golint at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

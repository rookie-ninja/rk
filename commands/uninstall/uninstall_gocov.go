package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install gocov on target hosts
func UninstallGoCovCommand() *cli.Command {
	command := &cli.Command{
		Name:      "gocov",
		Usage:     "uninstall gocov on local machine",
		UsageText: "rk uninstall gocov",
		Action:    UninstallGoCovAction,
	}

	return command
}

func UninstallGoCovAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-gocov")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("gocov", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall gocov at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

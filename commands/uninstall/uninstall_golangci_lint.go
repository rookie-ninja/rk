package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install golangci-lint on target hosts
func UninstallGoLangCILintCommand() *cli.Command {
	command := &cli.Command{
		Name:      "golangci-lint",
		Usage:     "uninstall golangci-lint on local machine",
		UsageText: "rk uninstall golangci-lint",
		Action:    UninstallGoLangCILintAction,
	}

	return command
}

func UninstallGoLangCILintAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-golangci-lint")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("golangci-lint", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall golangci-lint at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

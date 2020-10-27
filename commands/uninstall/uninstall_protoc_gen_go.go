package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install protoc-gen-go on target hosts
func UninstallProtocGenGoCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-go",
		Usage:     "uninstall protoc-gen-go on local machine",
		UsageText: "rk uninstall protoc-gen-go",
		Action:    UninstallProtocGenGoAction,
	}

	return command
}

func UninstallProtocGenGoAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-protoc-gen-go")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("protoc-gen-go", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall protoc-gen-go at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

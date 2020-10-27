package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install protoc-gen-swagger on target hosts
func UninstallProtocGenSwaggerCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-swagger",
		Usage:     "uninstall protoc-gen-swagger on local machine",
		UsageText: "rk uninstall protoc-gen-swagger",
		Action:    UninstallProtocGenSwaggerAction,
	}

	return command
}

func UninstallProtocGenSwaggerAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-protoc-gen-swagger")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("protoc-gen-swagger", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall protoc-gen-swagger at %s", path)
	exec.Command("rm", path).CombinedOutput()
	Success()
	return nil
}

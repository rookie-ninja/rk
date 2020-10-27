package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install protobuf on target hosts
func UninstallProtobufCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protobuf",
		Usage:     "uninstall protobuf on local machine",
		UsageText: "rk uninstall protobuf",
		Action:    UninstallProtobufAction,
	}

	return command
}

func UninstallProtobufAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("uninstall-protobuf")
	defer rk_common.Finish(event, nil)

	// check path
	path := CheckPath("protoc", event)

	if len(path) < 1 {
		Success()
		return nil
	}

	// uninstall release
	color.Cyan("Uninstall protobuf at %s", path)
	exec.Command("rm", path).CombinedOutput()
	exec.Command("rm", "-rf", "/usr/local/include/google")
	Success()
	return nil
}

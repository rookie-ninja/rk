// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package docker

import (
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
)

// Create docker run command
func dockerRun() *cli.Command {
	command := &cli.Command{
		Name:      "run",
		Usage:     "run a docker image built with rk docker build",
		UsageText: "rk docker run",
		Before:    common.CommandBefore,
		After:     common.CommandAfter,
		Action:    runAction,
	}

	return command
}

// Follow bellow sequence.
// 1: Compile source codes
// 2: Build docker image
// 3: Run docker image
func runAction(ctx *cli.Context) error {
	if err := common.UnmarshalBootConfig("build.yaml", common.BuildConfig); err != nil {
		return err
	}
	common.BuildConfig.Build.GOOS = "linux"
	common.BuildConfig.Build.GOARCH = "amd64"

	chain := common.NewActionChain()
	chain.Add("Validate docker environment", validateDockerCommand, false)
	chain.Add("Clearing target folder", func(ctx *cli.Context) error {
		return os.RemoveAll(common.BuildTarget)
	}, false)

	switch common.BuildConfig.Build.Type {
	case "go":
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Build go file", build.BuildGoFile, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	default:
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	}

	chain.Add("Build docker image", buildDockerImage, false)
	chain.Add("Run docker image", runDockerImage, false)

	return chain.Execute(ctx)

}

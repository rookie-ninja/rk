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

// Create docker build command
func dockerBuild() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "build a docker image built with rk build command",
		UsageText: "rk docker build",
		Before:    common.CommandBefore,
		After:     common.CommandAfter,
		Action:    buildAction,
	}

	return command
}

// Build codes first and then build docker image based on codes
func buildAction(ctx *cli.Context) error {
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

	// We can call actions in build command, however, in order to print sequence of actions,
	// we copied the same codes here.
	// TODO: Optimize bellow code
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

	return chain.Execute(ctx)
}

// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package docker

import (
	"errors"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
)

func dockerBuild() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "build a docker image built with rk build command",
		UsageText: "rk docker build",
		Before:    beforeDefault,
		After:     afterDefault,
		Action:    buildAction,
	}

	return command
}

func buildAction(ctx *cli.Context) error {
	if err := common.UnmarshalBootConfig("build.yaml", common.BuildConfig); err != nil {
		return err
	}

	common.BuildConfig.Build.GOOS = "linux"
	common.BuildConfig.Build.GOARCH = "amd64"

	chain := common.NewActionChain()
	chain.Add("Validate docker environment", validateDockerCommand, false)
	chain.Add("Clearing target folder", func(ctx *cli.Context) error {
		// 0: Not dir of where go.mod file exists
		if !rkcommon.FileExists("go.mod") {
			return errors.New("not a go directory, failed to lookup go.mod file")
		}
		return os.RemoveAll(common.BuildTarget)
	}, false)

	switch common.BuildConfig.Build.Type {
	case "go":
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Build go file", build.BuildGoFile, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Generate rk meta from on local", build.WriteRkMetaFile, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	default:
		chain.Add("Execute user command before", build.ExecCommandsBefore, false)
		chain.Add("Execute user script before", build.ExecScriptBefore, false)
		chain.Add("Copy to target folder", build.CopyToTarget, false)
		chain.Add("Generate rk meta from on local", build.WriteRkMetaFile, false)
		chain.Add("Execute user script after", build.ExecScriptAfter, false)
		chain.Add("Execute user command after", build.ExecCommandsAfter, false)
	}

	chain.Add("Build docker image", buildDockerImage, false)

	return chain.Execute(ctx)
}

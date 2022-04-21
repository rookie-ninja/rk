// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package build

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
)

// Build build current project based build.yaml file
func Build() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "Build project which contains build.yaml",
		UsageText: "rk build",
		Before:    common.CommandBefore,
		After:     common.CommandAfter,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug mod",
			},
		},
		Action: buildAction,
	}

	return command
}

func buildAction(ctx *cli.Context) error {
	// 1: read build.yaml file
	if err := common.UnmarshalBootConfig("build.yaml", common.BuildConfig); err != nil {
		return err
	}
	chain := common.NewActionChain()
	chain.Add("Clearing target folder", func(ctx *cli.Context) error {
		return os.RemoveAll(common.BuildTarget)
	}, false)

	// 2: Currently, support go only.
	// Just run user command, copy and script actions only.
	switch common.BuildConfig.Build.Type {
	case "go":
		chain.Add("Execute user command before", ExecCommandsBefore, false)
		chain.Add("Execute user script before", ExecScriptBefore, false)
		chain.Add("Build go file", BuildGoFile, false)
		chain.Add("Copy to target folder", CopyToTarget, false)
		chain.Add("Execute user script after", ExecScriptAfter, false)
		chain.Add("Execute user command after", ExecCommandsAfter, false)
	default:
		chain.Add("Execute user command before", ExecCommandsBefore, false)
		chain.Add("Execute user script before", ExecScriptBefore, false)
		chain.Add("Copy to target folder", CopyToTarget, false)
		chain.Add("Execute user script after", ExecScriptAfter, false)
		chain.Add("Execute user command after", ExecCommandsAfter, false)
	}

	return chain.Execute(ctx)
}

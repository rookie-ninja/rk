// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package build

import (
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
)

func Build() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "Build project which contains build.yaml",
		UsageText: "rk build",
		Before:    beforeDefault,
		After:     afterDefault,
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
		// 0: Move to dir of where go.mod file exists
		if err := os.Chdir(rkcommon.GetGoWd()); err != nil {
			return err
		}
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
		chain.Add("Generate rk meta from on local", WriteRkMetaFile, false)
		chain.Add("Execute user script after", ExecScriptAfter, false)
		chain.Add("Execute user command after", ExecCommandsAfter, false)
	default:
		chain.Add("Execute user command before", ExecCommandsBefore, false)
		chain.Add("Execute user script before", ExecScriptBefore, false)
		chain.Add("Copy to target folder", CopyToTarget, false)
		chain.Add("Generate rk meta from on local", WriteRkMetaFile, false)
		chain.Add("Execute user script after", ExecScriptAfter, false)
		chain.Add("Execute user command after", ExecCommandsAfter, false)
	}

	return chain.Execute(ctx)
}

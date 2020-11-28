// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_build

import (
	"github.com/fatih/color"
	rk_query "github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
)

type buildInfo struct {
	Debug bool
}

var BuildInfo = buildInfo{}

func Build() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "build project which contains build.yaml",
		UsageText: "rk build",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &BuildInfo.Debug,
			},
		},
		Action: BuildAction,
	}

	return command
}

func BuildAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("build")

	// 1: read build.yaml file
	config := &rk_common.BootConfig{}
	return DoBuildGo(config, event)
}

// a helper function
func DoBuildGo(config *rk_common.BootConfig, event rk_query.Event) error {
	color.Cyan("[Action] marshaling boot.yaml")
	if err := rk_common.MarshalBuildConfig("build.yaml", config, event); err != nil {
		event.AddPair("marshal-config", "fail")
		return rk_common.Error(event, err)
	}
	rk_common.Success()

	// 2: switch to find correct one
	switch config.Build.Type {
	case "go":
		if err := BuildGo(config, event); err != nil {
			return nil
		}
	default:
		if err := BuildDefault(config, event); err != nil {
			return nil
		}
	}

	return nil
}

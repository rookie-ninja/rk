// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_build

import "github.com/urfave/cli/v2"

type buildInfo struct {
	Debug bool
}

var BuildInfo = buildInfo{}

func Build() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "build project",
		UsageText: "rk build",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &BuildInfo.Debug,
			},
		},
		Subcommands: []*cli.Command{
			BuildGoCommand(),
		},
	}

	return command
}

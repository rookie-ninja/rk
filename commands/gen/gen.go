// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_gen

import "github.com/urfave/cli/v2"

func Gen() *cli.Command {
	command := &cli.Command{
		Name:      "gen",
		Usage:     "Generate files",
		UsageText: "rk gen [third-party file]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &GenInfo.Debug,
			},
		},

		Subcommands: []*cli.Command{
			GenPbGoCommand(),
			GenPbDocCommand(),
		},
	}

	return command
}

type genInfo struct {
	Debug bool
}

var GenInfo = genInfo{}

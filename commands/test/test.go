// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_test

import "github.com/urfave/cli/v2"

func Test() *cli.Command {
	command := &cli.Command{
		Name:      "test",
		Usage:     "Run unit test",
		UsageText: "rk test [third-party file]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &TestInfo.Debug,
			},
		},
		Subcommands: []*cli.Command{
			UtGoCommand(),
		},
	}

	return command
}

type testInfo struct {
	Debug bool
}

var TestInfo = testInfo{}

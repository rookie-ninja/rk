// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package docker

import (
	"github.com/urfave/cli/v2"
)

func Docker() *cli.Command {
	command := &cli.Command{
		Name:      "docker",
		Usage:     "Build a docker image built with rk build command",
		UsageText: "rk docker [build or run]",
		Subcommands: []*cli.Command{
			dockerBuild(),
			dockerRun(),
		},
	}

	return command
}

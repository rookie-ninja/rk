// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_docker

import (
	"github.com/urfave/cli/v2"
)

var (
	dockerRegistry = "rk-docker-image"
	dockerTag      = "latest"
)

type dockerInfo struct {
	Debug bool
}

var DockerInfo = dockerInfo{}

func Docker() *cli.Command {
	command := &cli.Command{
		Name:      "docker",
		Usage:     "build a docker image built with rk build command",
		UsageText: "rk docker [build or run]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &DockerInfo.Debug,
			},
		},
		Subcommands: []*cli.Command{
			DockerBuild(),
			DockerRun(),
		},
	}

	return command
}

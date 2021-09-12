// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

// Package pkg is responsible for creating rk style package.
package pkg

import (
	"github.com/urfave/cli/v2"
)

func Pkg() *cli.Command {
	command := &cli.Command{
		Name:      "pkg",
		Usage:     "Create, list or describe RK package",
		UsageText: "rk pkg [commands]",
		Subcommands: []*cli.Command{
			PkgList(),
			PkgDesc(),
			PkgCreate(),
		},
	}

	return command
}

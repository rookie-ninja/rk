// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/commands/clear"
	"github.com/rookie-ninja/rk/commands/docker"
	"github.com/rookie-ninja/rk/commands/install"
	"github.com/rookie-ninja/rk/commands/pack"
	"github.com/rookie-ninja/rk/commands/run"
	"github.com/rookie-ninja/rk/commands/uninstall"
	"github.com/rookie-ninja/rk/commands/ut"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "rk",
		Usage:       "rk utility command line tools",
		Version:     "1.2.4",
		Description: "rk is command line interface for utility tools during rk style software development lifecycle",
		Commands: []*cli.Command{
			build.Build(),
			clear.Clear(),
			docker.Docker(),
			install.Install(),
			uninstall.Uninstall(),
			pack.Pack(),
			run.Run(),
			ut.UT(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
		color.Yellow("Please read helper message by executing [rk -h]")
	}
}

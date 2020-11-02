// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/commands/gen"
	"github.com/rookie-ninja/rk/commands/install"
	"github.com/rookie-ninja/rk/commands/test"
	"github.com/rookie-ninja/rk/commands/uninstall"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "rk",
		Usage:       "rk utility command line tools",
		Version:     "1.0.0",
		Description: "rk is command line interface for utility tools during rk style software development lifecycle",
		Commands: []*cli.Command{
			rk_install.Install(),
			rk_gen.Gen(),
			rk_test.Test(),
			rk_build.Build(),
			rk_uninstall.Uninstall(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

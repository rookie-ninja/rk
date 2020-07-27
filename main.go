package main

import (
	"github.com/rookie-ninja/rk-cmd/commands/install"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name: "rk",
		Usage: "rk utility command line tools",
		Description: "rk is ad command line interface for utility tools during rk style software development lifecycle",
		Commands: []*cli.Command {
			rk_install.Install(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
package rk_gen

import "github.com/urfave/cli/v2"

func Gen() *cli.Command {
	command := &cli.Command{
		Name: "gen",
		Usage: "Generate files",
		UsageText: "rk gen [third-party file]",
		Flags: []cli.Flag {
			&cli.BoolFlag{
				Name: "debug",
				Aliases: []string{"d"},
				Usage: "debug mod",
				Destination: &GenInfo.Debug,
			},
		},

		Subcommands: []*cli.Command{
			GenPbGoCommand(),
		},
	}

	return command
}

type genInfo struct {
	Debug bool
}

var GenInfo = genInfo{}

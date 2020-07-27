package rk_install

import "github.com/urfave/cli/v2"

func Install() *cli.Command {
	command := &cli.Command{
		Name: "install",
		Aliases: []string{"i"},
		Usage: "install given third-party software",
		UsageText: "rk install [third-party software]",
		Subcommands: []*cli.Command{
			InstallProtobufCommand(),
			InstallProtocGenGoCommand(),
		},
	}

	return command
}

package rk_uninstall

import "github.com/urfave/cli/v2"

func Uninstall() *cli.Command {
	command := &cli.Command{
		Name:      "uninstall",
		Usage:     "Uninstall third-party software",
		UsageText: "rk uninstall [third-party software]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &InstallInfo.Debug,
			},
		},

		Subcommands: []*cli.Command{
			UninstallSwagCommand(),
			UninstallGoCovCommand(),
			UninstallGoLintCommand(),
			UninstallGoLangCILintCommand(),
			UninstallMockGenCommand(),
			UninstallProtobufCommand(),
			UninstallProtocGenDocCommand(),
			UninstallProtocGenGoCommand(),
			UninstallProtocGenSwaggerCommand(),
			UninstallRkStdCommand(),
		},
	}

	return command
}

type installInfo struct {
	Debug        bool
	ListReleases bool
	Release      string
}

var InstallInfo = installInfo{}

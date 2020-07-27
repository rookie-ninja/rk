package rk_install

import "github.com/urfave/cli/v2"

func Install() *cli.Command {
	command := &cli.Command{
		Name: "install",
		Usage: "Install third-party software",
		UsageText: "rk install [third-party software]",
		Flags: []cli.Flag {
			&cli.BoolFlag{
				Name: "debug",
				Aliases: []string{"d"},
				Usage: "debug mod",
				Destination: &InstallInfo.Debug,
			},
		},

		Subcommands: []*cli.Command{
			InstallProtobufCommand(),
			InstallProtocGenGoCommand(),
			InstallProtocGenGrpcGatewayCommand(),
			InstallProtocGenSwaggerCommand(),
			InstallGoLintCommand(),
			InstallGoLangCILintCommand(),
			InstallGoCovCommand(),
			InstallMockGenCommand(),
			InstallRkStdCommand(),
		},
	}

	return command
}

type installInfo struct {
	Debug bool
}

var InstallInfo = installInfo{}
// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import "github.com/urfave/cli/v2"

func Install() *cli.Command {
	command := &cli.Command{
		Name:      "install",
		Usage:     "Install third-party software",
		UsageText: "rk install [third-party software]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Aliases:     []string{"d"},
				Usage:       "debug mod",
				Destination: &InstallInfo.Debug,
			},
		},

		Subcommands: []*cli.Command{
			InstallProtobufCommand(),
			InstallProtocGenGoCommand(),
			InstallProtocGenGrpcGatewayCommand(),
			InstallProtocGenSwaggerCommand(),
			InstallProtocGenDocCommand(),
			InstallGoLintCommand(),
			InstallGoLangCILintCommand(),
			InstallGoCovCommand(),
			InstallMockGenCommand(),
			InstallSwagCommand(),
			InstallRkStdCommand(),
			InstallCfsslCommand(),
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

// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package uninstall

import "github.com/urfave/cli/v2"

func Uninstall() *cli.Command {
	command := &cli.Command{
		Name:      "uninstall",
		Usage:     "Uninstall third-party software",
		UsageText: "rk uninstall [third-party software]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug mod",
			},
		},

		Subcommands: []*cli.Command{
			uninstallBuf(),
			uninstallCfssl(),
			uninstallCfsslJson(),
			uninstallGoCov(),
			uninstallGolangCiLint(),
			uninstallMockGen(),
			uninstallPkger(),
			uninstallProtobuf(),
			uninstallProtocGenDoc(),
			uninstallProtocGenGo(),
			uninstallProtocGenGoGrpc(),
			uninstallProtocGenGrpcGateway(),
			uninstallProtocGenOpenApiV2(),
			uninstallSwag(),
			UninstallRkStdCommand(),
		},
	}

	return command
}

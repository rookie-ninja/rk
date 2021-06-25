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
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug mode",
			},
		},

		Subcommands: []*cli.Command{
			installProtobuf(),
			installProtocGenGo(),
			installProtocGenGrpcGateway(),
			installProtocGenOpenApiV2(),
			installProtocGenDoc(),
			installBuf(),
			installGolangCiLint(),
			installGoCov(),
			installMockGen(),
			installSwag(),
			installCfssl(),
			installCfsslJson(),
			installPkger(),
			installGolint(),
			rkStdCommand(),
		},
	}

	return command
}

// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package install

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
			installBuf(),
			installCfssl(),
			installCfsslJson(),
			installGoCov(),
			installGolangCiLint(),
			installMockGen(),
			installPkger(),
			installProtobuf(),
			installProtocGenDoc(),
			installProtocGenGo(),
			installProtocGenGrpcGateway(),
			installProtocGenOpenApiV2(),
			installSwag(),
			rkStdCommand(),
		},
	}

	return command
}

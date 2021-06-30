// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package ut

import (
	"fmt"
	"github.com/fatih/color"
	rkcommon "github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"path"
	"strings"
)

var (
	Lang string
)

func UT() *cli.Command {
	command := &cli.Command{
		Name:      "ut",
		Usage:     "Run unit test",
		UsageText: "rk ut",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lang, l",
				Aliases:     []string{"l"},
				Destination: &Lang,
				Required:    false,
				Usage:       "programming language, default:go",
			},
		},

		Action: utAction,
	}

	return command
}

func utAction(ctx *cli.Context) error {
	chain := common.NewActionChain()

	if len(Lang) < 1 {
		Lang = "go"
	}

	switch Lang {
	case "go":
		chain.Add("Running Go unit test", runGoUt, false)
		chain.Add(fmt.Sprintf("Generating unit test coverage report at %s", path.Base(rkcommon.RkUtHtmlFilePath)), runGoUtDoc, false)
	}

	return chain.Execute(ctx)
}

func runGoUt(ctx *cli.Context) error {
	args := make([]string, 0)
	args = append(args, "test", "./...", fmt.Sprintf("-coverprofile=%s", path.Base(rkcommon.RkUtOutFilepath)))

	color.White("- go %s", strings.Join(args, " "))

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.White(string(bytes))
		return err
	}

	color.White(string(bytes))

	return nil
}

func runGoUtDoc(ctx *cli.Context) error {
	args := make([]string, 0)
	args = append(args, "tool", "cover", fmt.Sprintf("-html=%s", path.Base(rkcommon.RkUtOutFilepath)), "-o", path.Base(rkcommon.RkUtHtmlFilePath))

	color.White("- go %s", strings.Join(args, " "))

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.White(string(bytes))
		return err
	}

	color.White(string(bytes))

	return nil
}

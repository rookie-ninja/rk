package ut

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
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
		chain.Add(fmt.Sprintf("Generating unit test coverage report at %s", "cov.html"), runGoUtDoc, false)
	}

	return chain.Execute(ctx)
}

func runGoUt(ctx *cli.Context) error {
	args := make([]string, 0)
	args = append(args, "test", "./...", fmt.Sprintf("-coverprofile=%s", "cov.out"))

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
	args = append(args, "tool", "cover", fmt.Sprintf("-html=%s", "cov.out"), "-o", "cov.html")

	color.White("- go %s", strings.Join(args, " "))

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		color.White(string(bytes))
		return err
	}

	color.White(string(bytes))

	return nil
}

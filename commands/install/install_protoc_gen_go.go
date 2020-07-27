package rk_install

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	rk_common "github.com/rookie-ninja/rk-cmd/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

//go get github.com/golang/protobuf/protoc-gen-go

// Install golang on target hosts
func InstallProtocGenGoCommand() *cli.Command{
	command := &cli.Command{
		Name: "protoc-gen-go",
		Usage: "protoc-gen-go",
		UsageText: "rk install protoc-gen-go",
		Subcommands: []*cli.Command{
			listProtocGenGoReleasesCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "version",
				Aliases: []string{"v"},
				Value: "",
				Usage: "install with version",
			},
		},
		Action: InstallProtocGenGoAction,
	}

	return command
}

func listProtocGenGoReleasesCommand() *cli.Command {
	command := &cli.Command{
		Name: "list",
		Usage: "list",
		UsageText: "rk install protoc-gen-go list",
		Action: listProtocGenGoReleasesAction,
	}

	return command
}

func listProtocGenGoReleasesAction(ctx *cli.Context) error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(),
		"golang", "protobuf", opt)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	for i := range res {
		buf.WriteString(res[i].GetTagName() + "\n")
	}

	color.Yellow(buf.String())

	return nil
}

func InstallProtocGenGoAction(ctx *cli.Context) error {
	var version = ""

	if version = ctx.String("version"); len(version) < 1 {
		color.Yellow("Install latest protoc-gen-go")
		version = "latest"
	}

	// install pkg
	color.Cyan("1: Install package")
	err := installProtocGenGo(version)
	if err != nil {
		color.Red("err: %v", err)
		return err
	}
	color.Green("Done")

	return nil
}

func installProtocGenGo(version string) error {
	url := "github.com/golang/protobuf/protoc-gen-go"

	if version != "latest" {
		url = fmt.Sprintf("github.com/golang/protobuf/protoc-gen-go@%s", version)
	}

	bytes, err := exec.Command("go", "get", url).CombinedOutput()
	if err != nil {
		return err
	}

	println(string(bytes))

	return nil
}
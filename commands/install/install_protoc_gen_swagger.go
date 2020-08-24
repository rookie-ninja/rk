package rk_install

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

const (
	ProtocGenSwaggerOwner   = "grpc-ecosystem"
	ProtocGenSwaggerRepo    = "grpc-gateway"
	ProtocGenSwaggerUrlBase = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
)

type installProtocGenSwaggerInfo struct {
	ListReleases bool
	Release      string
}

var InstallProtocGenSwaggerInfo = installProtocGenSwaggerInfo{
	ListReleases: false,
	Release:      "",
}

// Install protoc-gen-swagger on target hosts
func InstallProtocGenSwaggerCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-swagger",
		Usage:     "install protoc-gen-swagger on local machine",
		UsageText: "rk install protoc-gen-swagger -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallProtocGenSwaggerInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-swagger release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallProtocGenSwaggerInfo.ListReleases,
				Usage:       "list protoc-gen-swagger releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenSwaggerAction,
	}

	return command
}

func listProtocGenSwaggerReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtocGenSwaggerOwner, ProtocGenSwaggerRepo, opt)
	if err != nil {
		color.Red("failed to list protoc-gen-swagger release from github\n[err] %v", err)
		return err
	}

	for i := range res {
		color.Cyan("%s", res[i].GetTagName())
	}

	if len(res) < 1 {
		color.Cyan("no release found, please use default release")
	}

	return nil
}

func InstallProtocGenSwaggerAction(ctx *cli.Context) error {
	if InstallProtocGenSwaggerInfo.ListReleases {
		return listProtocGenSwaggerReleases()
	}

	// install pkg
	color.Cyan("Install protoc-gen-swagger %s", InstallProtocGenSwaggerInfo.Release)
	err := installProtocGenSwagger()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateProtocGenSwaggerInstallation()
}

func installProtocGenSwagger() error {
	url := fmt.Sprintf(ProtocGenSwaggerUrlBase)

	if len(InstallProtocGenSwaggerInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallProtocGenSwaggerInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install protoc-gen-swagger\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateProtocGenSwaggerInstallation() error {
	bytes, err := exec.Command("protoc-gen-swagger", "-version").CombinedOutput()
	if err != nil {
		color.Red("failed to validate protoc-gen-swagger\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

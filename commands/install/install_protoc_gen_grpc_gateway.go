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
	ProtocGenGrpcGatewayOwner   = "grpc-ecosystem"
	ProtocGenGrpcGatewayRepo    = "grpc-gateway"
	ProtocGenGrpcGatewayUrlBase = "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
)

type installProtocGenGrpcGatewayInfo struct {
	ListReleases bool
	Release      string
}

var InstallProtocGenGrpcGatewayInfo = installProtocGenGrpcGatewayInfo{
	ListReleases: false,
	Release:      "",
}

// Install protoc-gen-grpc-gateway on target hosts
func InstallProtocGenGrpcGatewayCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-grpc-gateway",
		Usage:     "install protoc-gen-grpc-gateway on local machine",
		UsageText: "rk install protoc-gen-grpc-gateway -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallProtocGenGrpcGatewayInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-grpc-gateway release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallProtocGenGrpcGatewayInfo.ListReleases,
				Usage:       "list protoc-gen-grpc-gateway releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenGrpcGatewayAction,
	}

	return command
}

func listProtocGenGrpcGatewayReleases() error {
	opt := &github.ListOptions{
		PerPage: 10,
	}

	res, _, err := rk_common.GithubClient.Repositories.ListReleases(context.Background(), ProtocGenGrpcGatewayOwner, ProtocGenGrpcGatewayRepo, opt)
	if err != nil {
		color.Red("failed to list protoc-gen-grpc-gateway release from github\n[err] %v", err)
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

func InstallProtocGenGrpcGatewayAction(ctx *cli.Context) error {
	if InstallProtocGenGrpcGatewayInfo.ListReleases {
		return listProtocGenGrpcGatewayReleases()
	}

	// install pkg
	color.Cyan("Install protoc-gen-grpc-gateway %s", InstallProtocGenGrpcGatewayInfo.Release)
	err := installProtocGenGrpcGateway()
	if err != nil {
		return err
	}
	color.Green("[Success]")

	color.Cyan("Validate installation")
	return validateProtocGenGrpcGatewayInstallation()
}

func installProtocGenGrpcGateway() error {
	url := fmt.Sprintf(ProtocGenGrpcGatewayUrlBase)

	if len(InstallProtocGenGrpcGatewayInfo.Release) > 1 {
		url += fmt.Sprintf("@%s", InstallProtocGenGrpcGatewayInfo.Release)
	}

	bytes, err := exec.Command("go", "get", "-v", url).CombinedOutput()
	if err != nil {
		color.Red("failed to install protoc-gen-grpc-gateway\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	if InstallInfo.Debug {
		color.Blue("[stdout] %s", string(bytes))
	}

	return nil
}

func validateProtocGenGrpcGatewayInstallation() error {
	bytes, err := exec.Command("protoc-gen-grpc-gateway", "-version").CombinedOutput()
	if err != nil {
		color.Red("failed to validate protoc-gen-grpc-gateway\n[err] %v", err)
		color.Red("[stderr] %s", string(bytes))
		return err
	}

	color.Green("[%s]", strings.TrimRight(string(bytes), "\n"))
	return nil
}

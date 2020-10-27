package rk_install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

const (
	ProtocGenGoOwner   = "golang"
	ProtocGenGoRepo    = "protobuf"
	ProtocGenGoUrlBase = "github.com/golang/protobuf/protoc-gen-go"
)

// Install protoc-gen-go on target hosts
func InstallProtocGenGoCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-go",
		Usage:     "install protoc-gen-go on local machine",
		UsageText: "rk install protoc-gen-go -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-go release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protoc-gen-go releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenGoAction,
	}

	return command
}

func InstallProtocGenGoAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protoc-gen-go-release")
		return PrintReleasesFromGithub(ProtocGenGoOwner, ProtocGenGoRepo, event)
	}

	event := rk_common.GetEvent("install-protoc-gen-go")
	// go get package
	if err := GoGetFromGithub(ProtocGenGoRepo, ProtocGenGoUrlBase, InstallInfo.Release, event); err != nil {
		return err
	}
	Success()

	if err := ValidateInstallation(exec.Command("which", "protoc-gen-go"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

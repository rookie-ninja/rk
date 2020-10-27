package rk_install

import (
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"runtime"
	"strings"
)

const (
	ProtocGenDocOwner = "pseudomuto"
	ProtocGenDocRepo  = "protoc-gen-doc"
)

// Install protoc-gen-grpc-gateway on target hosts
func InstallProtocGenDocCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protoc-gen-doc",
		Usage:     "install protoc-gen-doc on local machine",
		UsageText: "rk install protoc-gen-doc -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Required:    false,
				Usage:       "protoc-gen-doc release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protoc-gen-doc releases, list most recent 10 releases",
			},
		},
		Action: InstallProtocGenDocAction,
	}

	return command
}

func InstallProtocGenDocAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protoc-gen-doc-release")
		return PrintReleasesFromGithub(ProtocGenDocOwner, ProtocGenDocRepo, event)
	}

	event := rk_common.GetEvent("install-protoc-gen-doc")
	// 1: get release context
	releaseCtx, err := GetReleaseContext(ProtocGenDocOwner, ProtocGenDocRepo, InstallInfo.Release, event)
	if err != nil {
		return err
	}
	Success()

	// 2.1: specify filter based on pattern
	// https://github.com/pseudomuto/protoc-gen-doc/releases
	releaseCtx.Filter = func(url string) bool {
		var sysArch string
		if runtime.GOOS == "darwin" {
			sysArch = "darwin-amd64"
		} else if runtime.GOOS == "win" {
			sysArch = "windows-amd64"
		} else {
			sysArch = "linux-amd64"
		}

		if strings.Contains(url, sysArch) {
			return true
		}

		return false
	}

	// 2.2: download with filter
	DownloadGithubRelease(releaseCtx, event)
	if err != nil {
		return err
	}
	Success()

	// 3: install release
	color.Cyan("Install protoc-gen-go in %s to %s", releaseCtx.LocalFilePath, UserLocalBin)
	releaseCtx.ExtractType = "tar"
	releaseCtx.ExtractPath = "protoc-gen-doc"
	releaseCtx.ExtractArg = "--strip-components=1"
	if err = ExtractToDest(releaseCtx, event); err != nil {
		return err
	}
	Success()

	// 4: validate installation
	if err := ValidateInstallation(exec.Command("protoc-gen-doc", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

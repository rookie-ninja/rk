// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
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
	ProtobufOwner = "protocolbuffers"
	ProtobufRepo  = "protobuf"
)

// Install protobuf on target hosts
func InstallProtobufCommand() *cli.Command {
	command := &cli.Command{
		Name:      "protobuf",
		Usage:     "install protobuf on local machine",
		UsageText: "rk install protobuf -r [release]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "release, r",
				Aliases:     []string{"r"},
				Destination: &InstallInfo.Release,
				Usage:       "protobuf release",
			},
			&cli.BoolFlag{
				Name:        "list, l",
				Aliases:     []string{"l"},
				Destination: &InstallInfo.ListReleases,
				Usage:       "list protobuf releases, list most recent 10 releases",
			},
		},
		Action: InstallProtobufAction,
	}

	return command
}

func InstallProtobufAction(ctx *cli.Context) error {
	if InstallInfo.ListReleases {
		event := rk_common.GetEvent("list-protobuf-release")
		return PrintReleasesFromGithub(ProtobufOwner, ProtobufRepo, event)
	}

	event := rk_common.GetEvent("install-protobuf")
	// 1: get release context
	releaseCtx, err := GetReleaseContext(ProtobufOwner, ProtobufRepo, InstallInfo.Release, event)
	if err != nil {
		return err
	}
	Success()

	// 2: download from github
	// 2.1: specify filter based on pattern
	// https://github.com/protocolbuffers/protobuf/releases
	releaseCtx.Filter = func(url string) bool {
		var sysArch string
		if runtime.GOOS == "darwin" {
			sysArch = "osx-x86_64"
		} else if runtime.GOOS == "win" {
			sysArch = "win64"
		} else {
			sysArch = "linux-x86_64"
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
	color.Cyan("install protoc in %s to %s", releaseCtx.LocalFilePath, UserLocalBin)
	releaseCtx.ExtractType = "zip"
	releaseCtx.ExtractPath = "bin/protoc"
	if err = ExtractToDest(releaseCtx, event); err != nil {
		return err
	}
	releaseCtx.ExtractPath = "include/*"
	releaseCtx.DestPath = UserLocalInclude
	if err = ExtractToDest(releaseCtx, event); err != nil {
		return err
	}
	Success()

	// 4: validate installation
	if err := ValidateInstallation(exec.Command("protoc", "--version"), event); err != nil {
		return err
	}

	rk_common.Finish(event, nil)
	return nil
}

// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"runtime"
	"strings"
)

func installProtobuf() *cli.Command {
	command := commandDefault("protobuf")
	command.Before = beforeDefault
	command.Action = protobufAction
	command.After = afterDefault

	return command
}

func protobufAction(ctx *cli.Context) error {
	GithubInfo.Owner = "protocolbuffers"
	GithubInfo.Repo = "protobuf"
	GithubInfo.DecompressType = "zip"
	GithubInfo.DecompressedFilesToCopy = []string{
		"bin/protoc",
		"include/*",
	}
	GithubInfo.DecompressedFilesDestination = []string{
		UserLocalBin,
		UserLocalInclude,
	}
	GithubInfo.ValidationCmd = exec.Command("protoc", "--version")
	GithubInfo.DownloadFilter = func(url string) bool {
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

	// List tags only
	if hasListFlag(ctx) {
		chain := common.NewActionChain()
		chain.Add("List tags from github", printTagsFromGithub, false)
		return chain.Execute(ctx)
	}

	chain := common.NewActionChain()
	chain.Add("Find release from github", getReleaseToInstallFromGithub, false)
	chain.Add("Download from github", downloadFromGithub, false)
	chain.Add("Decompress file", decompressFile, false)
	chain.Add("Copy to destination", copyToDestination, false)
	chain.Add("Validate installation", validateInstallation, false)
	err := chain.Execute(ctx)

	// Log to event
	event := common.GetEvent(ctx)
	event.AddPayloads(githubInfoToPayloads()...)

	return err
}

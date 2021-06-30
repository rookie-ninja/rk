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

// Install on local machine
func installBuf() *cli.Command {
	command := commandDefault("buf")
	command.Before = beforeDefault
	command.Action = bufAction
	command.After = afterDefault

	return command
}

func bufAction(ctx *cli.Context) error {
	GithubInfo.Owner = "bufbuild"
	GithubInfo.Repo = "buf"
	GithubInfo.DecompressType = "tar"
	GithubInfo.DecompressedFilesToCopy = []string{
		"buf/bin/*",
	}
	GithubInfo.DecompressedFilesDestination = []string{
		UserLocalBin,
	}
	GithubInfo.ValidationCmd = exec.Command("buf", "--version")

	// List tags only
	if hasListFlag(ctx) {
		chain := common.NewActionChain()
		chain.Add("List tags from github", printTagsFromGithub, false)
		return chain.Execute(ctx)
	}

	GithubInfo.DownloadFilter = func(url string) bool {
		var sysArch string
		if runtime.GOOS == "darwin" {
			sysArch = "Darwin-x86_64.tar.gz"
		} else if runtime.GOOS == "win" {
			sysArch = ""
		} else {
			sysArch = "Linux-x86_64.tar.gz"
		}
		if strings.Contains(url, sysArch) {
			return true
		}

		return false
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

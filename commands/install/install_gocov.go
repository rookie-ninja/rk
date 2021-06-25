// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install on local machine
func installGoCov() *cli.Command {
	command := commandDefault("gocov")
	command.Before = beforeDefault
	command.Action = goCovAction
	command.After = afterDefault

	return command
}

func goCovAction(ctx *cli.Context) error {
	GithubInfo.Owner = "axw"
	GithubInfo.Repo = "gocov"
	GithubInfo.GoGetUrl = "github.com/axw/gocov/gocov"
	GithubInfo.ValidationCmd = exec.Command("which", "gocov")

	// List tags only
	if hasListFlag(ctx) {
		chain := rk_common.NewActionChain()
		chain.Add("List tags from github", printTagsFromGithub, false)
		return chain.Execute(ctx)
	}

	chain := rk_common.NewActionChain()
	chain.Add("Go get from remote repo", goGetFromRemoteUrl, false)
	chain.Add("Validate installation", validateInstallation, false)
	err := chain.Execute(ctx)

	// Log to event
	event := rk_common.GetEventV2(ctx)
	event.AddPayloads(githubInfoToPayloads()...)

	return err
}

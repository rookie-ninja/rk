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
func installCfssl() *cli.Command {
	command := commandDefault("cfssl")
	command.Before = beforeDefault
	command.Action = cfsslAction
	command.After = afterDefault

	return command
}

func cfsslAction(ctx *cli.Context) error {
	GithubInfo.Owner = "cloudflare"
	GithubInfo.Repo = "cfssl"
	GithubInfo.GoGetUrl = "github.com/cloudflare/cfssl/cmd/cfssl"
	GithubInfo.ValidationCmd = exec.Command("cfssl", "version")

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

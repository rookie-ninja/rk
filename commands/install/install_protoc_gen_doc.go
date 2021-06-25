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
func installProtocGenDoc() *cli.Command {
	command := commandDefault("protoc-gen-doc")
	command.Before = beforeDefault
	command.Action = protocGenDocAction
	command.After = afterDefault

	return command
}

func protocGenDocAction(ctx *cli.Context) error {
	GithubInfo.Owner = "pseudomuto"
	GithubInfo.Repo = "protoc-gen-doc"
	GithubInfo.GoGetUrl = "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc"
	GithubInfo.ValidationCmd = exec.Command("protoc-gen-doc", "-version")

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

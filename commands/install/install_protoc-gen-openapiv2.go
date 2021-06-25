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
func installProtocGenOpenApiV2() *cli.Command {
	command := commandDefault("protoc-gen-openapiv2")
	command.Before = beforeDefault
	command.Action = protocGenOpenApiV2Action
	command.After = afterDefault

	return command
}

func protocGenOpenApiV2Action(ctx *cli.Context) error {
	GithubInfo.Owner = "grpc-ecosystem"
	GithubInfo.Repo = "grpc-gateway"
	GithubInfo.GoGetUrl = "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	GithubInfo.ValidationCmd = exec.Command("protoc-gen-openapiv2", "--version")

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

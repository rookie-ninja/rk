// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install on local machine
func installProtocGenGrpcGateway() *cli.Command {
	command := commandDefault("protoc-gen-grpc-gateway")
	command.Before = beforeDefault
	command.Action = protocGenGrpcGatewayAction
	command.After = afterDefault

	return command
}

func protocGenGrpcGatewayAction(ctx *cli.Context) error {
	GithubInfo.Owner = "golang"
	GithubInfo.Repo = "protobuf"
	GithubInfo.GoGetUrl = "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	GithubInfo.ValidationCmd = exec.Command("protoc-gen-grpc-gateway", "--version")

	// List tags only
	if hasListFlag(ctx) {
		chain := common.NewActionChain()
		chain.Add("List tags from github", printTagsFromGithub, false)
		return chain.Execute(ctx)
	}

	chain := common.NewActionChain()
	chain.Add("Go get from remote repo", goGetFromRemoteUrl, false)
	chain.Add("Validate installation", validateInstallation, false)
	err := chain.Execute(ctx)

	// Log to event
	event := common.GetEvent(ctx)
	event.AddPayloads(githubInfoToPayloads()...)

	return err
}

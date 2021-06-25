package rk_install

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
)

// Install on local machine
func installPkger() *cli.Command {
	command := commandDefault("pkger")
	command.Before = beforeDefault
	command.Action = pkgerAction
	command.After = afterDefault

	return command
}

func pkgerAction(ctx *cli.Context) error {
	GithubInfo.Owner = "markbates"
	GithubInfo.Repo = "pkger"
	GithubInfo.GoGetUrl = "github.com/markbates/pkger/cmd/pkger"
	GithubInfo.ValidationCmd = exec.Command("which", "pkger")

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

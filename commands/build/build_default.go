package rk_build

import (
	"github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
)

func BuildDefault(config *rk_common.BootConfig, event rkquery.Event) error {
	event.AddPair("type", "default")
	// 1: run command before
	if err := runUserCommand(config.Build.Commands.Before, "before", event); err != nil {
		return err
	}

	// 2: run script before
	if err := runUserScripts(config.Build.Commands.Before, "before", event); err != nil {
		return err
	}

	// 3: build go file and move them to target folder
	if err := compileGoFile(config, event); err != nil {
		return err
	}

	// 4: run command after
	if err := runUserCommand(config.Build.Commands.After, "after", event); err != nil {
		return err
	}

	// 5: run script after
	if err := runUserScripts(config.Build.Commands.Before, "after", event); err != nil {
		return err
	}
	rk_common.Success()

	rk_common.Finish(event, nil)
	return nil
}

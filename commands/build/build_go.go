package rk_build

import (
	"github.com/fatih/color"
	rk_query "github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
)

func BuildGo(config *rk_common.BootConfig, event rk_query.Event) error {
	event.AddPair("type", "go")

	// prerequisite: clear target folder
	rk_common.ClearTargetFolder("target", event)

	// 2: validate go environment
	color.Cyan("[Action] validating go environment")
	if err := rk_common.ValidateCommand("go", event); err != nil {
		return err
	}
	rk_common.Success()

	// 3: run command before
	color.Cyan("[Action] run user command before")
	if err := runUserCommand(config.Build.Commands.Before, "before", event); err != nil {
		return err
	}
	rk_common.Success()

	// 4: run script before
	color.Cyan("[Action] run user script before")
	if err := runUserScripts(config.Build.Scripts.Before, "before", event); err != nil {
		return err
	}
	rk_common.Success()

	// 5: run swag command
	color.Cyan("[Action] run swag for swagger document")
	runSwagCommand(config, event)
	rk_common.Success()

	// 6: build go file and move them to target folder
	color.Cyan("[Action] build go file")
	if err := compileGoFile(config, event); err != nil {
		return err
	}
	rk_common.Success()

	// 7: copy files
	color.Cyan("[Action] copy files")
	if err := copyUserFolder(config, event); err != nil {
		return err
	}
	rk_common.Success()

	// 8: run command after
	color.Cyan("[Action] run user command after")
	if err := runUserCommand(config.Build.Commands.After, "after", event); err != nil {
		return err
	}
	rk_common.Success()

	// 9: run script after
	color.Cyan("[Action] run user script after")
	if err := runUserScripts(config.Build.Scripts.After, "after", event); err != nil {
		return err
	}
	rk_common.Success()

	rk_common.Finish(event, nil)
	return nil
}

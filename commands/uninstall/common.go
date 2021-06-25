// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_uninstall

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

var UninstallInfo = uninstallInfo{}

type uninstallInfo struct {
	app string
}

func commandDefault(name string) *cli.Command {
	return &cli.Command{
		Name:      name,
		Usage:     fmt.Sprintf("uninstall %s from local machine", name),
		UsageText: fmt.Sprintf("rk uninstall %s", name),
	}
}

func beforeDefault(ctx *cli.Context) error {
	name := strings.Join(strings.Split(ctx.Command.FullName(), " "), "/")
	event := rk_common.CreateEvent(name)

	// Inject event into context
	ctx.Context = context.WithValue(ctx.Context, rk_common.EventKey, event)

	return nil
}

func afterDefault(ctx *cli.Context) error {
	rk_common.Finish(rk_common.GetEventV2(ctx), nil)
	return nil
}

func checkPath(ctx *cli.Context) error {
	bytes, err := exec.Command("which", UninstallInfo.app).CombinedOutput()
	if err != nil {
		return err
	}

	path := strings.TrimSuffix(string(bytes), "\n")

	color.White(fmt.Sprintf("- uninstall %s at path:%s", UninstallInfo.app, path))
	rk_common.GetEventV2(ctx).AddPair("path", path)
	return nil
}

func validateUninstallation(ctx *cli.Context) error {
	bytes, _ := exec.Command("which", UninstallInfo.app).CombinedOutput()

	color.White(string(bytes))
	return nil
}

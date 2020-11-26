// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_run

import (
	"fmt"
	rk_build "github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Run() *cli.Command {
	command := &cli.Command{
		Name:      "run",
		Usage:     "run server build by rk build",
		UsageText: "rk run",
		Action:    RunAction,
	}

	return command
}

func RunAction(ctx *cli.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	event := rk_common.GetEvent("run")

	// 1: build project first
	event.StartTimer("build")
	if err := rk_build.BuildAction(ctx); err != nil {
		event.AddPair("build", "fail")
		event.AddErr(err)
		event.EndTimer("build")
		return rk_common.Error(event, err)
	}
	event.EndTimer("build")

	// 2: move to target folder and run binary
	os.Chdir("target")
	cmd := exec.Command("./bin/app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		event.AddPair("docker-run", "fail")
		rk_common.ClearTargetFolder("target", event)
		rk_common.Error(event,
			rk_common.NewScriptError(
				fmt.Sprintf("failed to run docker image \n[err] %v", err)))
		return err
	}

	go func() {
		<-sigs
		cmd.Process.Signal(syscall.SIGINT)
	}()

	cmd.Wait()

	event.AddPair("run", "success")
	rk_common.Success()

	return nil
}

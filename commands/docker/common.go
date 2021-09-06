// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package docker

import (
	"bufio"
	"context"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func beforeDefault(ctx *cli.Context) error {
	name := strings.Join(strings.Split(ctx.Command.FullName(), " "), "/")
	event := common.CreateEvent(name)
	event.AddPayloads(zap.Strings("flags", ctx.FlagNames()))

	// Inject event into context
	ctx.Context = context.WithValue(ctx.Context, common.EventKey, event)

	return nil
}

func afterDefault(ctx *cli.Context) error {
	common.Finish(common.GetEvent(ctx), nil)
	return nil
}

func runDockerImage(ctx *cli.Context) error {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	args := []string{"run", "--rm"}
	if len(common.BuildConfig.Docker.Run.Args) > 0 {
		for i := range common.BuildConfig.Docker.Run.Args {
			arg := common.BuildConfig.Docker.Run.Args[i]
			if len(arg) > 0 {
				args = append(args, strings.Split(arg, " ")...)
			}
		}
	}

	tag := common.BuildConfig.Docker.Build.Registry + ":" + common.BuildConfig.Docker.Build.Tag

	args = append(args, tag)

	color.White("- docker %s", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		os.RemoveAll(common.BuildTarget)
		return err
	}

	go func() {
		<-sig
		cmd.Process.Signal(syscall.SIGINT)
	}()

	cmd.Wait()

	return nil
}

func buildDockerImage(ctx *cli.Context) error {
	meta := common.GetRkMetaFromCmd()

	if len(common.BuildConfig.Docker.Build.Registry) < 1 {
		common.BuildConfig.Docker.Build.Registry = meta.Name
	}
	if len(common.BuildConfig.Docker.Build.Tag) < 1 {
		common.BuildConfig.Docker.Build.Tag = meta.Version
	}

	args := []string{"build"}
	if len(common.BuildConfig.Docker.Build.Args) > 0 {
		for i := range common.BuildConfig.Docker.Build.Args {
			arg := common.BuildConfig.Docker.Build.Args[i]
			if len(arg) > 0 {
				args = append(args, strings.Split(arg, " ")...)
			}
		}
	}

	tag := common.BuildConfig.Docker.Build.Registry + ":" + common.BuildConfig.Docker.Build.Tag

	args = append(args, "-t", tag, ".")

	color.White("- docker %s", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)
	// out, _ := cmd.StdoutPipe()
	errOut, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		os.RemoveAll(common.BuildTarget)
		return err
	}

	go func() {
		scanner := bufio.NewScanner(errOut)
		for scanner.Scan() {
			// We received an error, assign error
			msg := scanner.Text()
			color.White(msg)
		}
	}()

	cmd.Wait()

	return nil
}

func validateDockerCommand(ctx *cli.Context) error {
	bytes, err := exec.Command("which", "docker").CombinedOutput()
	if err != nil {
		return err
	}

	color.Yellow(string(bytes))

	return nil
}

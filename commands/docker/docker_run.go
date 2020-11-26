package rk_docker

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func DockerRun() *cli.Command {
	command := &cli.Command{
		Name:      "run",
		Usage:     "run a docker image built with rk docker build",
		UsageText: "run",
		Action:    DockerRunAction,
	}

	return command
}

func DockerRunAction(ctx *cli.Context) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	event := rk_common.GetEvent("docker-run")

	// 1: read build.yaml file
	config := &rk_common.BootConfig{}
	if err := rk_common.MarshalBuildConfig("build.yaml", config, event); err != nil {
		event.AddPair("marshal-config", "fail")
		return rk_common.Error(event, err)
	}

	// 1: build project first
	event.StartTimer("docker-build")
	if err := DockerBuildAction(ctx); err != nil {
		event.AddPair("docker-build", "fail")
		event.AddErr(err)
		event.EndTimer("docker-build")
		return rk_common.Error(event, err)
	}
	event.EndTimer("docker-build")

	// 3: determine docker registry and tag
	if len(config.Docker.Build.Registry) < 1 {
		config.Docker.Build.Registry = dockerRegistry
	}
	if len(config.Docker.Build.Tag) < 1 {
		config.Docker.Build.Tag = dockerTag
	}

	args := []string{"run", "--rm"}
	if len(config.Docker.Run.Args) > 0 {
		for i := range config.Docker.Run.Args {
			arg := config.Docker.Run.Args[i]
			if len(arg) > 0 {
				args = append(args, strings.Split(arg, " ")...)
			}
		}
	}

	args = append(args, config.Docker.Build.Registry+":"+config.Docker.Build.Tag)

	color.Cyan(fmt.Sprintf("docker %s", strings.Join(args, " ")))

	cmd := exec.Command("docker", args...)
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

	rk_common.Success()
	event.AddPair("docker-run", "success")

	rk_common.Finish(event, nil)
	return nil
}

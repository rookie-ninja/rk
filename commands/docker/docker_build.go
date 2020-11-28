package rk_docker

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk/commands/build"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

func DockerBuild() *cli.Command {
	command := &cli.Command{
		Name:      "build",
		Usage:     "build a docker image built with rk build command",
		UsageText: "build",
		Action:    DockerBuildAction,
	}

	return command
}

func DockerBuildAction(ctx *cli.Context) error {
	event := rk_common.GetEvent("docker-build")

	// 0: validate docker command
	color.Cyan("[Action] validating docker environment")
	if err := rk_common.ValidateCommand("docker", event); err != nil {
		return err
	}
	rk_common.Success()

	// 1: read build.yaml file
	config := &rk_common.BootConfig{}
	if err := rk_common.MarshalBuildConfig("build.yaml", config, event); err != nil {
		event.AddPair("marshal-config", "fail")
		return rk_common.Error(event, err)
	}

	// 1: build project first
	event.StartTimer("build")
	config.Build.GOOS = "linux"
	config.Build.GOARCH = "amd64"
	if err := rk_build.DoBuildGo(config, event); err != nil {
		event.AddPair("build", "fail")
		event.AddErr(err)
		event.EndTimer("build")
		return rk_common.Error(event, err)
	}
	event.EndTimer("build")

	// 3: determine docker registry and tag
	if len(config.Docker.Build.Registry) < 1 {
		config.Docker.Build.Registry = dockerRegistry
	}
	if len(config.Docker.Build.Tag) < 1 {
		config.Docker.Build.Tag = dockerTag
	}

	args := []string{"build"}
	if len(config.Docker.Build.Args) > 0 {
		for i := range config.Docker.Build.Args {
			arg := config.Docker.Build.Args[i]
			if len(arg) > 0 {
				args = append(args, strings.Split(arg, " ")...)
			}
		}
	}

	args = append(args, "-t", config.Docker.Build.Registry+":"+config.Docker.Build.Tag, ".")

	color.Cyan(fmt.Sprintf("docker %s", strings.Join(args, " ")))

	cmd := exec.Command("docker", args...)
	out, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		event.AddPair("docker-build", "fail")
		rk_common.ClearTargetFolder("target", event)
		rk_common.Error(event,
			rk_common.NewScriptError(
				fmt.Sprintf("failed to build docker image \n[err] %v", err)))
		return err
	}

	go func() {
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			color.Yellow(scanner.Text())
		}
	}()

	cmd.Wait()

	rk_common.Success()
	event.AddPair("docker-build", "success")

	rk_common.Finish(event, nil)
	return nil
}

package rk_build

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	RkHomeDir     = "./"
	targetFolders = []string{
		"target",
		"target/bin",
		"target/cmd",
		"target/logs",
		"target/script",
		"target/config",
	}
)

func init() {
	userHomeDir, _ := os.UserHomeDir()
	RkHomeDir = path.Join(userHomeDir, ".rk")
	os.MkdirAll(RkHomeDir, os.ModePerm)
}

// run user commands
func runUserCommand(commands []string, pos string, event rk_query.Event) error {
	if len(commands) < 1 {
		color.Yellow("no user command")
		return nil
	}

	for i := range commands {
		cmd := commands[i]
		tokens := strings.Split(cmd, " ")
		if len(tokens) < 1 || len(tokens[0]) < 1 {
			color.Yellow(fmt.Sprintf("command-%s-%d is empty, skip", pos, i))
			continue
		}
		color.Yellow(fmt.Sprintf("command-%s-%d:[%s]", pos, i, cmd))
		event.AddFields(zap.Strings(fmt.Sprintf("cmd-%s-%d", pos, i), tokens))

		if len(tokens) > 1 {
			bytes, err := exec.Command(tokens[0], tokens[1:]...).CombinedOutput()
			if err != nil {
				rk_common.ClearTargetFolder("target", event)
				wrap := rk_common.NewCommandError(
					fmt.Sprintf("failed to run command %s, clearing...\n[err] %v\n[stderr] %s", cmd, err, string(bytes)))
				return rk_common.Error(event, wrap)
			}
			color.White(string(bytes))
		} else {
			bytes, err := exec.Command(tokens[0]).CombinedOutput()

			if err != nil {
				rk_common.ClearTargetFolder("target", event)
				wrap := rk_common.NewCommandError(
					fmt.Sprintf("failed to run command %s, clearing... \n[err] %v\n[stderr] %s", cmd, err, string(bytes)))
				return rk_common.Error(event, wrap)
			}
			color.White(string(bytes))
		}
	}

	event.AddPair("cmd-"+pos, "success")
	return nil
}

// run user scripts
func runUserScripts(scripts []string, pos string, event rk_query.Event) error {
	if len(scripts) < 1 {
		color.Yellow("no user script")
		return nil
	}

	event.AddFields(zap.Strings("scripts", scripts))

	for i := range scripts {
		script := scripts[i]

		if len(script) < 1 {
			color.Yellow(fmt.Sprintf("script-%s-%d is empty, skip", pos, i))
			continue
		}

		color.Yellow(fmt.Sprintf("script-%s-%d:[%s]", pos, i, script))
		bytes, err := exec.Command("sh", script).CombinedOutput()
		if err != nil {
			rk_common.ClearTargetFolder("target", event)
			return rk_common.Error(event,
				rk_common.NewScriptError(
					fmt.Sprintf("failed to run script %v, clearing... \n[err] %v\n[stderr] %s", script, err, string(bytes))))
		}
		color.White(string(bytes))
	}

	event.AddPair("script-"+pos, "success")
	return nil
}

// run swag for swagger documentation
func runSwagCommand(config *rk_common.BootConfig, event rk_query.Event) {
	if len(config.Build.Main) < 1 {
		// Main function path is empty, let's assuming main function located at
		// current directory
		config.Build.Main = "."
	}

	if err := rk_common.ValidateCommand("swag", event); err != nil {
		// just skip current stage since swag is missing
		color.Yellow(err.Error())
		event.AddErr(err)
		event.AddPair("swag", "skip")
		return
	}

	bytes, err := exec.Command("swag", "init", "-g", config.Build.Main).CombinedOutput()
	if err != nil {
		color.Yellow(string(bytes))
		color.Yellow(err.Error())
		event.AddErr(err)
		event.AddPair("swag", "skip")
	}
	event.AddPair("swag", "success")
}

// run go build command
func compileGoFile(config *rk_common.BootConfig, event rk_query.Event) error {
	// 1: create directory named as target and sub folders
	for i := range targetFolders {
		if err := os.MkdirAll(targetFolders[i], os.ModePerm); err != nil {
			rk_common.ClearTargetFolder("target", event)
			return rk_common.Error(event,
				rk_common.NewScriptError(
					fmt.Sprintf("failed to create %s directory, clearing... \n[err] %v", targetFolders[i], err)))
		}
	}

	// 2: locate source main.go file
	if len(config.Build.Main) < 1 {
		// Main function path is empty, let's assuming main function located at
		// current directory
		config.Build.Main = "main.go"
	}

	// 3: determine GOOS
	if len(config.Build.GOOS) > 0 {
		os.Setenv("GOOS", config.Build.GOOS)
	}

	// 4: determine GOARCH
	if len(config.Build.GOARCH) > 0 {
		os.Setenv("GOARCH", config.Build.GOARCH)
	}

	// 5: run go build command with args and move it to target/bin folder
	args := []string{
		"build",
		"-o",
		"target/bin/app",
	}

	// append user provided args
	for i := range config.Build.Args {
		arg := config.Build.Args[i]
		if len(arg) > 0 {
			args = append(args, strings.Split(arg, " ")...)
		}
	}

	// append main.go files
	args = append(args, config.Build.Main)

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		rk_common.ClearTargetFolder("target", event)
		return rk_common.Error(event,
			rk_common.NewScriptError(
				fmt.Sprintf("failed to build go binary to target folder, clearing... \n[err] %v \n[stderr] %v", err, string(bytes))))
	}

	// 6: copy boot.yaml file to target folder
	if err := copyFolder("boot.yaml", "target", event); err != nil {
		rk_common.ClearTargetFolder("target", event)
		return rk_common.Error(event,
			rk_common.NewScriptError(
				fmt.Sprintf("failed to copy boot.yaml to target folder \n[err] %v \n[stderr] %v", err, string(bytes))))
	}

	// 7: copy docs to target folder
	if err := copyFolder("docs", "target", event); err != nil {
		rk_common.ClearTargetFolder("target", event)
		return rk_common.Error(event,
			rk_common.NewScriptError(
				fmt.Sprintf("failed to copy docs to target folder \n[err] %v \n[stderr] %v", err, string(bytes))))
	}

	return nil
}

// copy user specified folder to target folder
func copyUserFolder(config *rk_common.BootConfig, event rk_query.Event) error {
	if len(config.Build.Copy) < 1 {
		color.Yellow("no folders or files need to copy, skip")
		return nil
	}

	event.AddFields(zap.Strings("copy", config.Build.Copy))

	for i := range config.Build.Copy {
		if len(config.Build.Copy[i]) < 1 {
			color.Yellow(fmt.Sprintf("copy-%d is empty, skip", i))
			continue
		}

		color.Yellow(fmt.Sprintf("copy-%d:[%s]", i, config.Build.Copy[i]))

		if err := copyFolder(config.Build.Copy[i], "target/", event); err != nil {
			rk_common.ClearTargetFolder("target", event)
			return rk_common.Error(event,
				rk_common.NewScriptError(
					fmt.Sprintf("failed to copy %s, clearing... \n[err] %v", config.Build.Copy[i], err)))
		}
	}

	return nil
}

// helper function for cp command
func copyFolder(src string, dst string, event rk_query.Event) error {
	// validate existence of src folder
	if _, err := os.Stat(src); err != nil {
		return err
	}

	bytes, err := exec.Command("cp", "-r", src, dst).CombinedOutput()
	if err != nil {
		event.AddPair(fmt.Sprintf("copy-%s", src), "fail")
		color.Red(string(bytes))
		return err
	}

	event.AddPair(fmt.Sprintf("copy-%s", src), "success")

	return nil
}

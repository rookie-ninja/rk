// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	targetFolders = []string{
		"target",
		"target/bin",
	}
)

// ExecCommandsBefore Execute user provided commands before building
func ExecCommandsBefore(ctx *cli.Context) error {
	return execCommands(common.BuildConfig.Build.Commands.Before)
}

// ExecCommandsAfter Execute user provided commands after building
func ExecCommandsAfter(ctx *cli.Context) error {
	return execCommands(common.BuildConfig.Build.Commands.After)
}

// Helper function to exec commands
func execCommands(commands []string) error {
	if len(commands) < 1 {
		color.White("- No user commands found!")
		return nil
	}

	for i := range commands {
		cmd := commands[i]
		tokens := strings.Split(cmd, " ")
		if len(tokens) < 1 || len(tokens[0]) < 1 {
			continue
		}

		color.White(fmt.Sprintf("- %s", cmd))

		if len(tokens) > 1 {
			bytes, err := exec.Command(tokens[0], tokens[1:]...).CombinedOutput()
			if err != nil {
				os.RemoveAll(common.BuildTarget)
				return err
			}
			color.White(string(bytes))
		} else {
			bytes, err := exec.Command(tokens[0]).CombinedOutput()
			if err != nil {
				os.RemoveAll(common.BuildTarget)
				return err
			}
			color.White(string(bytes))
		}
	}

	return nil
}

// ExecScriptBefore Execute user provided scripts before building
func ExecScriptBefore(ctx *cli.Context) error {
	return execScripts(common.BuildConfig.Build.Scripts.Before)
}

// ExecScriptAfter Execute user provided scripts after building
func ExecScriptAfter(ctx *cli.Context) error {
	return execScripts(common.BuildConfig.Build.Scripts.After)
}

func execScripts(scripts []string) error {
	if len(scripts) < 1 {
		color.White("- No user scripts found!")
		return nil
	}

	for i := range scripts {
		script := scripts[i]

		if len(script) < 1 {
			continue
		}

		color.White(fmt.Sprintf("- %s", script))
		bytes, err := exec.Command("sh", script).CombinedOutput()
		if err != nil {
			os.RemoveAll(common.BuildTarget)
			return err
		}
		color.White(string(bytes))
	}

	return nil
}

// WriteRkMetaFile Write RK meta into .rk folder
func WriteRkMetaFile(ctx *cli.Context) error {
	meta := common.GetRkMetaFromCmd()

	// 5: Write to .rk/git.yaml file
	color.White("- Write files to %s", path.Join(common.BuildTarget, rkcommon.RkMetaFilePath))
	if raw, err := yaml.Marshal(meta); err != nil {
		color.Yellow("[Warn] %v", err)
	} else {
		if err := ioutil.WriteFile(path.Join(common.BuildTarget, rkcommon.RkMetaFilePath), raw, os.ModePerm); err != nil {
			color.Yellow("[Warn] %v", err)
		}
	}

	return nil
}

// BuildGoFile Build go file
func BuildGoFile(ctx *cli.Context) error {
	// Create .rk folder first
	os.MkdirAll(path.Join(common.BuildTarget, path.Dir(rkcommon.RkMetaFilePath)), os.ModePerm)

	// 0: Not dir of where go.mod file exists
	if !rkcommon.FileExists("go.mod") {
		return errors.New("not a go directory, failed to lookup go.mod file")
	}

	// 1: create directory named as target and sub folders
	for i := range targetFolders {
		if err := os.MkdirAll(targetFolders[i], os.ModePerm); err != nil {
			os.RemoveAll(common.BuildTarget)
			return nil
		}
	}

	// 2: locate source main.go file
	if len(common.BuildConfig.Build.Main) < 1 {
		// Main function path is empty, let's assuming main function located at current directory
		common.BuildConfig.Build.Main = "main.go"
	}

	// 3: determine GOOS
	if len(common.BuildConfig.Build.GOOS) > 0 {
		os.Setenv("GOOS", common.BuildConfig.Build.GOOS)
	}

	// 4: determine GOARCH
	if len(common.BuildConfig.Build.GOARCH) > 0 {
		os.Setenv("GOARCH", common.BuildConfig.Build.GOARCH)
	}

	// 5: run go build command with args and move it to target/bin folder
	args := []string{
		"build",
		"-o",
		fmt.Sprintf("target/bin/%s", common.TargetPkgName),
	}

	// append user provided args
	for i := range common.BuildConfig.Build.Args {
		arg := common.BuildConfig.Build.Args[i]
		if len(arg) > 0 {
			args = append(args, strings.Split(arg, " ")...)
		}
	}

	// append main.go files
	args = append(args, common.BuildConfig.Build.Main)

	color.White("- go %s", strings.Join(args, " "))

	bytes, err := exec.Command("go", args...).CombinedOutput()
	if len(bytes) > 0 {
		color.White(string(bytes))
	}

	if err != nil {
		os.RemoveAll(common.BuildTarget)
		return err
	}

	// 6: copy boot.yaml file to target folder
	copyCommand("boot.yaml", common.BuildTarget)

	// 7: copy README.md, go.mod, LICENSE and unit test coverage report file .rk dir if exist
	copyCommand("README.md", path.Join(common.BuildTarget, rkcommon.RkReadmeFilePath))
	copyCommand("LICENSE", path.Join(common.BuildTarget, rkcommon.RkLicenseFilePath))
	copyCommand("go.mod", path.Join(common.BuildTarget, rkcommon.RkDepFilePath))
	copyCommand(path.Base(rkcommon.RkUtHtmlFilePath), path.Join(common.BuildTarget, rkcommon.RkUtHtmlFilePath))

	return nil
}

// CopyToTarget Copy folders to target folder
func CopyToTarget(ctx *cli.Context) error {
	if len(common.BuildConfig.Build.Copy) < 1 {
		color.White("- No folders need to copy!")
		return nil
	}

	for i := range common.BuildConfig.Build.Copy {
		folder := common.BuildConfig.Build.Copy[i]

		if len(folder) < 1 {
			continue
		}

		if err := copyCommand(folder, common.BuildTarget+"/"); err != nil {
			os.RemoveAll(common.BuildTarget)
			return err
		}
	}

	return nil
}

// CopyCommand Copy folders to target folder
func copyCommand(src string, dst string) error {
	os.MkdirAll(path.Dir(dst), os.ModePerm)

	args := []string{
		"-r", src, dst,
	}

	color.White("- cp %s", strings.Join(args, " "))
	bytes, err := exec.Command("cp", args...).CombinedOutput()
	if len(bytes) > 0 {
		color.White(string(bytes))
	}
	return err
}

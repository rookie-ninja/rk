// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk/common"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path"
)

const (
	pkgerIdentity           = "github.com/rookie-ninja/rk"
	pkgerRootDir            = "/commands/pkg/assets/"
	memFsRootDir            = "standard/"
	remotePkgRepoFromGithub = "https://github.com/rookie-ninja/rk-demo"
	remotePkgRepoFromGitee  = "https://gitee.com/rookie-ninja/rk-demo"
)

// PkgCreate Create package
func PkgCreate() *cli.Command {
	command := &cli.Command{
		Name:      "create",
		Usage:     "Create package",
		UsageText: "rk pkg create -t [type] -n [name] -d [destination] -rs [from remote source]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "package type hope to create, default is go-gin",
				Value:   "go-gin",
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "package name, default is demo",
				Value:   "demo",
			},
			&cli.StringFlag{
				Name:    "destination",
				Aliases: []string{"d"},
				Usage:   "where to create on local machine, default is ./",
				Value:   ".",
			},
			&cli.BoolFlag{
				Name:    "remote-source",
				Aliases: []string{"rs"},
				Usage:   "force to pull newest package from remote source, default is false",
				Value:   false,
			},
		},
		Action: pkgCreateAction,
	}

	return command
}

// Action of list sub command
func pkgCreateAction(ctx *cli.Context) error {
	fromRemoteSource := getPkgFromRemoteSourceFromFlags(ctx)

	chain := common.NewActionChain()
	if fromRemoteSource {
		chain.Add("Create package from remote source", copyFromRemoteToLocal, false)
	} else {
		chain.Add("Create package from local cache", copyFromCacheToLocal, false)
	}
	chain.Add("Validating local environment and provide suggestion", suggestion, true)

	return chain.Execute(ctx)
}

func copyFromRemoteToLocal(ctx *cli.Context) error {
	pkgType := getPkgTypeFromFlags(ctx)
	pkgName := getPkgNameFromFlags(ctx)
	pkgDest := getPkgDestinationFromFlags(ctx)

	// Copy from remote git repo to mem fs
	opts := &git.CloneOptions{
		URL:      remotePkgRepoFromGithub,
		Progress: os.Stdout,
	}

	srcFs := memfs.New()
	if _, err := git.Clone(memory.NewStorage(), srcFs, opts); err != nil {
		return err
	}

	c := rkcommon.NewMemFsToLocalFsCopier(srcFs)
	if err := c.CopyDir(path.Join(memFsRootDir, pkgType), path.Join(pkgDest, pkgName)); err != nil {
		return err
	}

	return nil
}

func copyFromCacheToLocal(ctx *cli.Context) error {
	pkgType := getPkgTypeFromFlags(ctx)
	pkgName := getPkgNameFromFlags(ctx)
	pkgDest := getPkgDestinationFromFlags(ctx)

	c := rkcommon.NewPkgerFsToLocalFsCopier(pkgerIdentity)
	if err := c.CopyDir(path.Join(pkgerRootDir, pkgType), path.Join(pkgDest, pkgName)); err != nil {
		return err
	}

	return nil
}

// Validate tools on local machine and give suggestion to user.
// Add suggestion based on user selection
// If user select gin framework, then we will suggest bellow to be installed:
// - golang
// - docker
// - doctoc
// - golangci-lint
// - rk
// - swag
//
// If user select grpc framework, then we will suggest bellow to be installed:
// - golang
// - docker
// - doctoc
// - golangci-lint
// - rk
// - protoc-gen-go
// - protoc-gen-go-grpc
// - protoc-gen-grpc-gateway
// - protoc-gen-openapiv2
// - buf
func suggestion(ctx *cli.Context) error {
	// Validate common tools
	// golang
	checkToolOnLocalMachine("go", "Please visit https://golang.org/ to install it.")
	// docker
	checkToolOnLocalMachine("docker", "Please visit https://www.docker.com/ to install it.")
	// doctoc
	checkToolOnLocalMachine("doctoc", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://github.com/thlorenz/doctoc.")
	// golangci-lint
	checkToolOnLocalMachine("golangci-lint", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://github.com/golangci/golangci-lint.")
	// rk
	checkToolOnLocalMachine("rk", "Please visit https://github.com/rookie-ninja/rk.")

	// 1: Get framework flag
	if getPkgTypeFromFlags(ctx) == "gin" {
		// swag
		checkToolOnLocalMachine("swag", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://github.com/swaggo/swag.")
	} else {
		// protoc-gen-go
		checkToolOnLocalMachine("protoc-gen-go", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://grpc.io/docs/languages/go/quickstart/.")
		// protoc-gen-go-grpc
		checkToolOnLocalMachine("protoc-gen-go-grpc", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://grpc.io/docs/languages/go/quickstart/.")
		// protoc-gen-grpc-gateway
		checkToolOnLocalMachine("protoc-gen-grpc-gateway", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://github.com/grpc-ecosystem/grpc-gateway.")
		// protoc-gen-openapiv2
		checkToolOnLocalMachine("protoc-gen-openapiv2", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://grpc.io/docs/languages/go/quickstart/.")
		// buf
		checkToolOnLocalMachine("buf", "Please use [rk cli](https://github.com/rookie-ninja/rk) to install it or visit https://docs.buf.build/installation.")
	}

	return nil
}

// Validate suggested tools exist on local machine
func checkToolOnLocalMachine(tool, suggestion string) {
	if _, err := exec.Command("which", tool).CombinedOutput(); err != nil {
		// Failed to exec command, show related message
		color.Yellow(fmt.Sprintf("- [Not found] %s", tool))
		color.Yellow(fmt.Sprintf("[Suggestion] %s", suggestion))

		return
	}

	// Command found
	color.White(fmt.Sprintf("[Found] %s", tool))
}

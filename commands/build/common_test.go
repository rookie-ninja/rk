// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package build

import (
	"context"
	"flag"
	rkquery "github.com/rookie-ninja/rk-query"
	"github.com/rookie-ninja/rk/common"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestBeforeDefault(t *testing.T) {
	ctx := cli.NewContext(nil, flag.NewFlagSet("", flag.ContinueOnError), nil)
	ctx.Context = context.TODO()

	assert.Nil(t, beforeDefault(ctx))
	assert.NotNil(t, common.GetEvent(ctx))
}

func TestAfterDefault(t *testing.T) {
	ctx := cli.NewContext(nil, flag.NewFlagSet("", flag.ContinueOnError), nil)
	ctx.Context = context.TODO()

	beforeDefault(ctx)
	assert.Nil(t, afterDefault(ctx))
	assert.Equal(t, rkquery.Ended, common.GetEvent(ctx).GetEventStatus())
}

func TestExecCommandsBefore(t *testing.T) {
	// Happy case
	defer func() {
		common.BuildConfig.Build.Commands.Before = []string{}
	}()
	common.BuildConfig.Build.Commands.Before = []string{
		"echo ut",
	}
	assert.Nil(t, ExecCommandsBefore(nil))

	// Without commands
	common.BuildConfig.Build.Commands.Before = []string{}
	assert.Nil(t, ExecCommandsBefore(nil))

	// Without command args
	common.BuildConfig.Build.Commands.Before = []string{
		"echo",
	}
	assert.Nil(t, ExecCommandsBefore(nil))

	// With invalid command
	common.BuildConfig.Build.Commands.Before = []string{
		"invalid ut",
	}
	assert.NotNil(t, ExecCommandsBefore(nil))
	common.BuildConfig.Build.Commands.Before = []string{
		"invalid",
	}
	assert.NotNil(t, ExecCommandsBefore(nil))
}

func TestExecCommandsAfter(t *testing.T) {
	// Happy case
	defer func() {
		common.BuildConfig.Build.Commands.After = []string{}
	}()
	common.BuildConfig.Build.Commands.After = []string{
		"echo ut",
	}
	assert.Nil(t, ExecCommandsAfter(nil))

	// Without commands
	common.BuildConfig.Build.Commands.After = []string{}
	assert.Nil(t, ExecCommandsAfter(nil))

	// Without command args
	common.BuildConfig.Build.Commands.After = []string{
		"echo",
	}
	assert.Nil(t, ExecCommandsAfter(nil))

	// With invalid command
	common.BuildConfig.Build.Commands.After = []string{
		"invalid ut",
	}
	assert.NotNil(t, ExecCommandsAfter(nil))
	common.BuildConfig.Build.Commands.After = []string{
		"invalid",
	}
	assert.NotNil(t, ExecCommandsAfter(nil))
}

func TestExecScriptBefore(t *testing.T) {
	// With invalid script
	defer func() {
		common.BuildConfig.Build.Scripts.After = []string{}
	}()
	common.BuildConfig.Build.Scripts.Before = []string{
		"invalid-script.sh",
	}
	assert.NotNil(t, ExecScriptBefore(nil))

	// Happy case
	// Create a script and run it
	script := `
#!/bin/sh
echo "ut"
`
	defer func() {
		os.RemoveAll("ut.sh")
	}()
	assert.Nil(t, ioutil.WriteFile("ut.sh", []byte(script), os.ModePerm))
	common.BuildConfig.Build.Scripts.Before = []string{
		"ut.sh",
	}
	assert.Nil(t, ExecScriptBefore(nil))
}

func TestExecScriptAfter(t *testing.T) {
	// With invalid script
	defer func() {
		common.BuildConfig.Build.Scripts.After = []string{}
	}()
	common.BuildConfig.Build.Scripts.After = []string{
		"invalid-script.sh",
	}
	assert.NotNil(t, ExecScriptAfter(nil))

	// Happy case
	// Create a script and run it
	script := `
#!/bin/sh
echo "ut"
`
	defer func() {
		os.RemoveAll("ut.sh")
	}()
	assert.Nil(t, ioutil.WriteFile("ut.sh", []byte(script), os.ModePerm))
	common.BuildConfig.Build.Scripts.After = []string{
		"ut.sh",
	}
	assert.Nil(t, ExecScriptAfter(nil))
}

func TestWriteRkMetaFile(t *testing.T) {
	assert.Nil(t, WriteRkMetaFile(nil))
}

func TestCopyToTarget(t *testing.T) {
	defer func() {
		common.BuildConfig.Build.Copy = []string{}
	}()
	// Without copy
	assert.Nil(t, CopyToTarget(nil))

	// Happy case
	// Create target folder first
	assert.Nil(t, os.MkdirAll("target", os.ModePerm))

	assert.Nil(t, os.MkdirAll("ut", os.ModePerm))
	defer func() {
		os.RemoveAll("ut")
		os.RemoveAll("target")
	}()

	common.BuildConfig.Build.Copy = []string{
		"ut",
	}
	assert.Nil(t, CopyToTarget(nil))
}

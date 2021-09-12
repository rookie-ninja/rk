// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package common

import (
	"context"
	"errors"
	"github.com/rookie-ninja/rk-query"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestGetEvent(t *testing.T) {
	ctx := cli.NewContext(nil, nil, nil)
	event := rkquery.NewEventFactory().CreateEventNoop()
	ctx.Context = context.WithValue(context.TODO(), EventKey, event)

	// With Event in context
	assert.Equal(t, event, GetEvent(ctx))

	// With no event in context
	ctx = cli.NewContext(nil, nil, nil)
	ctx.Context = context.TODO()
	assert.NotNil(t, GetEvent(ctx))
}

func TestCreateEvent(t *testing.T) {
	event := CreateEvent("ut-operation")
	assert.NotNil(t, event)
	assert.Equal(t, "ut-operation", event.GetOperation())
	assert.Equal(t, rkquery.InProgress, event.GetEventStatus())
}

func TestError(t *testing.T) {
	defer assertNotPanic(t)

	Error(nil)
	Error(errors.New("ut-error"))
}

func TestFinish(t *testing.T) {
	// With nil error
	event := rkquery.NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())

	Finish(event, nil)
	assert.Equal(t, "OK", event.GetResCode())
	assert.Equal(t, rkquery.Ended, event.GetEventStatus())

	// With error
	event = rkquery.NewEventFactory().CreateEvent()
	event.SetStartTime(time.Now())
	Finish(event, errors.New("ut-error"))
	assert.Equal(t, "Failed", event.GetResCode())
	assert.Equal(t, rkquery.Ended, event.GetEventStatus())
}

func TestUnmarshalBootConfig(t *testing.T) {
	defer assertNotPanic(t)

	configYaml := `
---
build:
  type: go
`
	tmpDir := t.TempDir()
	assert.Nil(t, ioutil.WriteFile(path.Join(tmpDir, "build.yaml"), []byte(configYaml), os.ModePerm))

	// With empty config path
	assert.NotNil(t, UnmarshalBootConfig("", &BootConfigForBuild{}))

	// With nil config
	assert.NotNil(t, UnmarshalBootConfig(path.Join(tmpDir, "build.yaml"), nil))

	// With invalid config type
	assert.NotNil(t, UnmarshalBootConfig(path.Join(tmpDir, "build.yaml"), "invalid"))

	// Happy case
	config := &BootConfigForBuild{}
	assert.Nil(t, UnmarshalBootConfig(path.Join(tmpDir, "build.yaml"), config))
	assert.Equal(t, "go", config.Build.Type)
}

func TestGetGoPathBin(t *testing.T) {
	assert.Equal(t, "bin", path.Base(GetGoPathBin()))
	assert.NotEmpty(t, path.Base(GetGoPathBin()))
}

func TestGetPkgName(t *testing.T) {
	defer assertNotPanic(t)
	assert.NotEmpty(t, getPkgName())
}

func TestGetGitInfo(t *testing.T) {
	defer assertNotPanic(t)
	assert.NotNil(t, getGitInfo())
}

func TestGetRkMetaFromCmd(t *testing.T) {
	defer assertNotPanic(t)

	assert.NotNil(t, GetRkMetaFromCmd())
}

func assertNotPanic(t *testing.T) {
	if r := recover(); r != nil {
		// Expect panic to be called with non nil error
		assert.True(t, false)
	} else {
		// This should never be called in case of a bug
		assert.True(t, true)
	}
}

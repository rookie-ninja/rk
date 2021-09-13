// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package docker

import (
	"github.com/rookie-ninja/rk/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// It may not work on the CI docker image since it doesn't have docker daemon background.
// Make sure UT pass on local machine.
func TestBuildDockerImage(t *testing.T) {
	defer assertNotPanic(t)
	defer func() {
		common.BuildConfig.Docker.Build.Args = []string{}
		common.BuildConfig.Docker.Build.Registry = ""
		common.BuildConfig.Docker.Build.Tag = ""
		os.RemoveAll("Dockerfile")
	}()

	// Without dockerfile
	// Do not assert anything, since it may success on local machine but fail on CI image
	buildDockerImage(nil)

	// With dockerfile
	dockerfile := `
FROM alpine
`
	wd, _ := os.Getwd()
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "Dockerfile"), []byte(dockerfile), os.ModePerm))
	assert.NotEmpty(t, common.BuildConfig.Docker.Build.Registry)
	assert.NotEmpty(t, common.BuildConfig.Docker.Build.Tag)
	buildDockerImage(nil)

	// With args
	common.BuildConfig.Docker.Build.Args = []string{
		"-q",
	}
	assert.NotEmpty(t, common.BuildConfig.Docker.Build.Registry)
	assert.NotEmpty(t, common.BuildConfig.Docker.Build.Tag)
	buildDockerImage(nil)
}

func TestRunDockerImage(t *testing.T) {
	defer assertNotPanic(t)
	defer func() {
		common.BuildConfig.Docker.Run.Args = []string{}
		common.BuildConfig.Docker.Build.Registry = ""
		common.BuildConfig.Docker.Build.Tag = ""
		os.RemoveAll("Dockerfile")
	}()

	// Without docker image
	// Do not assert anything, since it may success on local machine but fail on CI image
	runDockerImage(nil)

	// With docker image provided
	// Do not assert anything, since it may success on local machine but fail on CI image
	common.BuildConfig.Docker.Build.Registry = "alpine"
	common.BuildConfig.Docker.Build.Tag = "latest"
	runDockerImage(nil)

	// With docker args
	// Do not assert anything, since it may success on local machine but fail on CI image
	common.BuildConfig.Docker.Build.Registry = "alpine"
	common.BuildConfig.Docker.Build.Tag = "latest"
	common.BuildConfig.Docker.Run.Args = []string{
		"-w",
		"/usr/",
	}
	runDockerImage(nil)
}

func TestValidateDockerCommand(t *testing.T) {
	defer assertNotPanic(t)
	// It may success on local machine with docker installed and fail on CI machine
	validateDockerCommand(nil)
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

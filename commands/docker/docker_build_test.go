// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package docker

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestBuildAction_WithOtherType(t *testing.T) {
	defer assertNotPanic(t)
	defer func() {
		os.RemoveAll("build.yaml")
		os.RemoveAll("Dockerfile")
	}()

	// Create fake build.yaml file
	buildYaml := `
---
build:
  type: other
`
	// With dockerfile
	dockerfile := `
FROM alpine
`
	wd, _ := os.Getwd()
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "build.yaml"), []byte(buildYaml), os.ModePerm))
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "Dockerfile"), []byte(dockerfile), os.ModePerm))
	ctx := cli.NewContext(nil, nil, nil)

	assert.Nil(t, buildAction(ctx))
}

func TestBuildAction_WithGo(t *testing.T) {
	defer func() {
		os.RemoveAll("build.yaml")
		os.RemoveAll("go.mod")
		os.RemoveAll("main.go")
		os.RemoveAll("target")
		os.RemoveAll("Dockerfile")
	}()

	wd, _ := os.Getwd()

	// Create fake build.yaml and go.mod file
	buildYaml := `
---
build:
  type: go
  copy: [""]
`
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "build.yaml"), []byte(buildYaml), os.ModePerm))

	goMod := `
module ut-test

go 1.16
`
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "go.mod"), []byte(goMod), os.ModePerm))

	mainGo := `
package main

func main() {}
`
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "main.go"), []byte(mainGo), os.ModePerm))

	// With dockerfile
	dockerfile := `
FROM alpine
`
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "Dockerfile"), []byte(dockerfile), os.ModePerm))

	ctx := cli.NewContext(nil, nil, nil)
	assert.Nil(t, buildAction(ctx))
}

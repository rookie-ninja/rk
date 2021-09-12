// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package build

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestBuild(t *testing.T) {
	assert.NotNil(t, Build())
}

func TestBuildAction_WithOtherType(t *testing.T) {
	defer func() {
		os.RemoveAll("build.yaml")
	}()

	// Create fake build.yaml file
	buildYaml := `
---
build:
  type: other
  copy: [""]
  commands:
    before:
      - "echo before"
    after:
      - "echo after"
`
	wd, _ := os.Getwd()
	assert.Nil(t, ioutil.WriteFile(path.Join(wd, "build.yaml"), []byte(buildYaml), os.ModePerm))
	ctx := cli.NewContext(nil, nil, nil)
	assert.Nil(t, buildAction(ctx))
}

func TestBuildAction_WithGo(t *testing.T) {
	defer func() {
		os.RemoveAll("build.yaml")
		os.RemoveAll("go.mod")
		os.RemoveAll("main.go")
		os.RemoveAll("target")
	}()

	wd, _ := os.Getwd()

	// Create fake build.yaml and go.mod file
	buildYaml := `
---
build:
  type: go
  copy: [""]
  commands:
    before:
      - "echo before"
    after:
      - "echo after"
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

	ctx := cli.NewContext(nil, nil, nil)
	assert.Nil(t, buildAction(ctx))
}

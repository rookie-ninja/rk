// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_uninstall

import (
	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/rookie-ninja/rk-query"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	UserLocalBin     = "/usr/local/bin"
	UserLocalInclude = "/usr/local/include"
)

var (
	GithubClient = github.NewClient(nil)
	RkHomeDir    = "./"
)

type URLFilter func(string) bool

type GithubReleaseContext struct {
	Release       *github.RepositoryRelease
	Repo          string
	Owner         string
	Filter        func(string) bool
	RemoteURL     string
	LocalFilePath string
	ExtractPath   string
	ExtractType   string
	ExtractArg    string
	TempPath      string
	DestPath      string
	AssetSize     int64
}

func init() {
	userHomeDir, _ := os.UserHomeDir()
	RkHomeDir = path.Join(userHomeDir, ".rk")
	os.MkdirAll(RkHomeDir, os.ModePerm)
}

func Success() {
	color.Green("[success]")
	color.White("--------------------------------")
}

func CheckPath(app string, event rk_query.Event) string {
	bytes, _ := exec.Command("which", app).CombinedOutput()

	path := strings.TrimSuffix(string(bytes), "\n")

	color.Cyan("Check %s path:%s", app, path)
	event.AddPair("path", path)
	return path
}

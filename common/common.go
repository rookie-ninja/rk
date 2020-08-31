package rk_common

import (
	"github.com/google/go-github/v32/github"
	"os"
	"path"
	"runtime"
)

var (
	GithubClient = github.NewClient(nil)
	RkHomeDir    = "./"
)

func init() {
	userHomeDir, _ := os.UserHomeDir()
	RkHomeDir = path.Join(userHomeDir, ".rk")
	os.MkdirAll(RkHomeDir, os.ModePerm)
}

func StdOut(input string) {
	println(input)
}

func GetSysAndArch() string {
	var sysArch string

	if runtime.GOOS == "darwin" {
		sysArch = "osx-x86_64"
	} else if runtime.GOOS == "win" {
		sysArch = "win64"
	} else {
		sysArch = "linux-x86_64"
	}
	return sysArch
}


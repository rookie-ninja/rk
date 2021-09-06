// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/modfile"
	rkcommon "github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go/build"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	EventKey    = "rkEvent"
	BuildTarget = "target"
)

var (
	confStr = `{
     "level": "info",
     "encoding": "console",
     "outputPaths": ["%s"],
     "errorOutputPaths": ["stderr"],
     "initialFields": {},
     "encoderConfig": {
       "messageKey": "msg",
       "levelKey": "",
       "nameKey": "",
       "timeKey": "",
       "callerKey": "",
       "stacktraceKey": "",
       "callstackKey": "",
       "errorKey": "",
       "timeEncoder": "iso8601",
       "fileKey": "",
       "levelEncoder": "capital",
       "durationEncoder": "second",
       "callerEncoder": "full",
       "nameEncoder": "full"
     },
    "maxsize": 1,
    "maxage": 7,
    "maxbackups": 3,
    "localtime": true,
    "compress": true
   }`

	logger        *zap.Logger
	factory       *rkquery.EventFactory
	BuildConfig   = &BootConfig{}
	TargetPkgName = getPkgName()
	TargetGitInfo = getGitInfo()
)

type BootConfig struct {
	Build struct {
		Type     string   `yaml:"type"`
		Main     string   `yaml:"main"`
		GOOS     string   `yaml:"GOOS"`
		GOARCH   string   `yaml:"GOARCH"`
		Args     []string `yaml:"args"`
		Copy     []string `yaml:"copy"`
		Commands struct {
			Before []string `yaml:"before"`
			After  []string `yaml:"after"`
		} `yaml:"commands"`
		Scripts struct {
			Before []string `yaml:"before"`
			After  []string `yaml:"after"`
		} `yaml:"scripts"`
	} `yaml:"build"`
	Docker struct {
		Build struct {
			Registry string   `yaml:"registry"`
			Tag      string   `yaml:"tag"`
			Args     []string `yaml:"args"`
		} `yaml:"build"`
		Run struct {
			Args []string `yaml:"args"`
		} `yaml:"run"`
	} `yaml:"docker"`
}

func init() {
	userHomeDir, _ := os.UserHomeDir()
	rkDir := path.Join(userHomeDir, ".rk/logs")
	if err := os.MkdirAll(rkDir, os.ModePerm); err != nil {
		color.Yellow("Failed to create rk temp dir.")
		color.Red(err.Error())
	}

	bytes := []byte(fmt.Sprintf(confStr, path.Join(rkDir, "rk.log")))
	logger, _, _ = rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)
	factory = rkquery.NewEventFactory(
		rkquery.WithAppName("rk-cmd"),
		rkquery.WithZapLogger(logger))
}

func GetEvent(ctx *cli.Context) rkquery.Event {
	if v := ctx.Context.Value(EventKey); v != nil {
		if event, ok := v.(rkquery.Event); ok {
			return event
		}
	}

	return factory.CreateEventNoop()
}

func CreateEvent(op string) rkquery.Event {
	event := factory.CreateEvent()
	event.SetOperation(op)
	event.SetStartTime(time.Now())
	return event
}

func Finish(event rkquery.Event, err error) {
	if err != nil {
		event.AddErr(err)
		event.SetResCode("Failed")
	} else if len(event.GetResCode()) < 1 {
		event.SetResCode("OK")
	}

	event.SetEndTime(time.Now())
	event.Finish()
}

func Error(err error) error {
	if err != nil {
		color.Red("[Error] %s", err.Error())
	} else {
		color.Red("[Error] internal error occur")
	}

	return err
}

func ClearTargetFolder(path string, event rkquery.Event) {
	// just try our best to remove it
	if err := os.RemoveAll(path); err != nil {
		color.Red(err.Error())
		event.AddPair("clear-target", "fail")
	}

	event.AddPair("clear-target", "success")
}

func ClearDir(path string) {
	// just try our best to remove it
	if err := os.RemoveAll(path); err != nil {
		color.Red(err.Error())
	}
}

func UnmarshalBootConfig(configFilePath string, config interface{}) error {
	if !path.IsAbs(configFilePath) {
		// Ignore error
		wd, _ := os.Getwd()
		configFilePath = path.Join(wd, configFilePath)
	}

	bytes := TryReadFile(configFilePath)
	if len(bytes) < 1 {
		return errors.New(fmt.Sprintf("Failed to read build.yaml file from path:%s", configFilePath))
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		return err
	}

	return nil
}

func GetGoPathBin() string {
	return path.Join(build.Default.GOPATH, "bin")
}

// TryReadFile reads files with provided path, use working directory if given path is relative path.
// Ignoring error while reading.
func TryReadFile(filePath string) []byte {
	if len(filePath) < 1 {
		return make([]byte, 0)
	}

	if !path.IsAbs(filePath) {
		// Ignore error
		wd, _ := os.Getwd()
		filePath = path.Join(wd, filePath)
	}

	if bytes, err := ioutil.ReadFile(filePath); err != nil {
		return bytes
	} else {
		return bytes
	}
}

// getPkgName will try best to retrieve package name by bellow priority.
// 1: Read from go.mod file
// 2: Read from current directory
// 3: Return empty if nothing matched
func getPkgName() string {
	// Read from go.mod
	if rkcommon.FileExists("go.mod") {
		if bytes := rkcommon.TryReadFile("go.mod"); len(bytes) > 0 {
			return path.Base(modfile.ModulePath(bytes))
		}
	}

	// Read current directory name
	if res, err := os.Getwd(); err == nil {
		return path.Base(res)
	}

	// Return empty
	return ""
}

func getGitInfo() *rkcommon.Git {
	res := &rkcommon.Git{
		Commit: &rkcommon.Commit{
			Committer: &rkcommon.Committer{},
		},
	}

	opt := &git.PlainOpenOptions{
		DetectDotGit: true,
	}
	if repo, err := git.PlainOpenWithOptions(".", opt); err == nil {
		// Parse URL
		if config, e1 := repo.Config(); e1 == nil {
			// Iterate map
			for _, v := range config.Remotes {
				if len(v.URLs) > 0 {
					res.Url = strings.TrimSuffix(v.URLs[0], ".git")
				}
				break
			}
		}

		// Get current head
		var head plumbing.Hash
		if ref, e1 := repo.Head(); e1 == nil {
			head = ref.Hash()
		}

		// Parse branch
		if iter, e1 := repo.Branches(); e1 == nil {
			iter.ForEach(func(ref *plumbing.Reference) error {
				if ref.Hash() == head {
					res.Branch = ref.Name().Short()
				}
				return nil
			})
		}

		// Parse tag
		if tag, e1 := repo.TagObject(head); e1 == nil {
			res.Tag = tag.Name
		}

		// Parse commit
		if commit, e1 := repo.CommitObject(head); e1 == nil {
			res.Commit.Committer.Name = commit.Committer.Name
			res.Commit.Committer.Email = commit.Committer.Email
			res.Commit.Date = commit.Author.When.Format(object.DateFormat)
			res.Commit.Id = commit.ID().String()
			if len(res.Commit.Id) > 7 {
				res.Commit.IdAbbr = commit.ID().String()[:7]
			}
			res.Commit.Sub = commit.Message
		}
	}

	return res
}

// GetRkMetaFromCmd construct RkMeta from local environment
func GetRkMetaFromCmd() *rkcommon.RkMeta {
	meta := &rkcommon.RkMeta{}

	// Load package name from go.mod file
	meta.Name = TargetPkgName

	// Load current git info from local git
	meta.Git = TargetGitInfo

	// Load version
	if len(meta.Git.Tag) < 1 {
		// Tag is empty, let's try to construct version as <branch>-<commit IdAbbr>
		if len(meta.Git.Branch) < 1 {
			// Branch is empty, we will use commit IdAbbr instead
			meta.Version = meta.Git.Commit.IdAbbr
		} else {
			// <branch>-<commit IdAbbr>
			meta.Version = strings.Join([]string{meta.Git.Branch, meta.Git.Commit.IdAbbr}, "-")
		}
	} else {
		// Use tag as version
		meta.Version = meta.Git.Tag
	}

	return meta
}

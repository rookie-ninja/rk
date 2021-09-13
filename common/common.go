// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package common

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/modfile"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go/build"
	"os"
	"path"
	"strings"
	"time"
)

const (
	// Event key in context
	EventKey = "rkEvent"
	// The folder where build result will be moved to
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

	logger  *zap.Logger
	factory *rkquery.EventFactory
	// BuildConfig match yaml config of build.yaml
	BuildConfig = &BootConfigForBuild{}
	// TargetPkgName the compiled target name
	TargetPkgName = getPkgName()
	// TargetGitInfo git information while building
	TargetGitInfo = getGitInfo()
)

// BootConfigForBuild match yaml config of build.yaml
// 1: Build.Type: Required, the type of programming language
// 2: Build.Main: Optional, path of main.go file, default is current working directory
// 3: Build.GOOS: Optional, specify GOOS, default from current go env
// 4: Build.GOARCH: Optional, specify GOARCH, default from current go env
// 5: Build.Args: Optional, specify arguments while run build
// 6: Build.Copy: Optional, the files and folders need to copy to target folder
// 7: Build.Commands.Before: Optional, commands needs to run before build
// 8: Build.Commands.After: Optional, commands needs to run after build
// 9: Build.Scripts.Before: Optional, scripts needs to run before build
// 10: Build.Scripts.After: Optional, scripts needs to run after build
// 11: Docker.Build.Registry: Optional, docker registry while building docker image
// 12: Docker.Build.Tag: Optional, docker tag while building docker image
// 13: Docker.Build.Args: Optional, arguments needs to pass while building docker image
// 14: Docker.Build.Run.Args: Optional, arguments needs to pass to docker while running docker container
type BootConfigForBuild struct {
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
	// Get user home directory
	userHomeDir, _ := os.UserHomeDir()
	// rk cli home directory which is ~/.rk
	rkDir := path.Join(userHomeDir, ".rk/logs")
	// Make directory of ~/.rk/logs
	if err := os.MkdirAll(rkDir, os.ModePerm); err != nil {
		color.Yellow("Failed to create rk temp dir.")
		color.Red(err.Error())
	}
	// ZapLogger config which replace the output folder with ~/.rk/logs
	bytes := []byte(fmt.Sprintf(confStr, path.Join(rkDir, "rk.log")))
	// ZapLogger
	logger, _, _ = rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)
	// Create EventFactory
	factory = rkquery.NewEventFactory(
		rkquery.WithAppName("rk-cmd"),
		rkquery.WithZapLogger(logger))
}

// GetEvent Retrieve event in context
func GetEvent(ctx *cli.Context) rkquery.Event {
	if v := ctx.Context.Value(EventKey); v != nil {
		if event, ok := v.(rkquery.Event); ok {
			return event
		}
	}

	return factory.CreateEventNoop()
}

// CreateEvent Create an event with operation name
func CreateEvent(op string) rkquery.Event {
	event := factory.CreateEvent()
	event.SetOperation(op)
	event.SetStartTime(time.Now())
	return event
}

// Finish current event
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

// Error Print error with red color of text
func Error(err error) error {
	if err != nil {
		color.Red("[Error] %s", err.Error())
	} else {
		color.Red("[Error] internal error occur")
	}

	return err
}

// UnmarshalBootConfig Read build.yaml and unmarshal it into config
func UnmarshalBootConfig(configFilePath string, config interface{}) error {
	if config == nil {
		return errors.New("Nil config provided")
	}

	if !path.IsAbs(configFilePath) {
		// Ignore error
		wd, _ := os.Getwd()
		configFilePath = path.Join(wd, configFilePath)
	}

	bytes := rkcommon.TryReadFile(configFilePath)
	if len(bytes) < 1 {
		return errors.New(fmt.Sprintf("Failed to read build.yaml file from path:%s", configFilePath))
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		return err
	}

	return nil
}

// GetGoPathBin Get $GOPATH/bin path on local machine
func GetGoPathBin() string {
	return path.Join(build.Default.GOPATH, "bin")
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

// Get git information from .git folder if current package is a git package
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

// Create an event and inject into context
func CommandBefore(ctx *cli.Context) error {
	name := strings.Join(strings.Split(ctx.Command.FullName(), " "), "/")
	event := CreateEvent(name)
	event.AddPayloads(zap.Strings("flags", ctx.FlagNames()))

	// Inject event into context
	ctx.Context = context.WithValue(ctx.Context, EventKey, event)

	return nil
}

// Extract event and finish event
func CommandAfter(ctx *cli.Context) error {
	Finish(GetEvent(ctx), nil)
	return nil
}

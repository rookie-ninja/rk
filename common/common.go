// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	rkcommon "github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"path"
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

	logger      *zap.Logger
	factory     *rkquery.EventFactory
	BuildConfig = &BootConfig{}
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

	bytes := rkcommon.TryReadFile(configFilePath)
	if len(bytes) < 1 {
		return errors.New(fmt.Sprintf("Failed to read build.yaml file from path:%s", configFilePath))
	}

	if err := yaml.Unmarshal(bytes, config); err != nil {
		return err
	}

	return nil
}

func GetGoPathBin() string {
	return path.Join(rkcommon.GetGoEnv().GoPath, "bin")
}
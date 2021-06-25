// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const (
	EventKey = "rkEvent"
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
	Pack struct {
		Name string `yaml:"name"`
	} `yaml:"pack"`
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
	//userHomeDir, _ := os.UserHomeDir()
	//rkDir := path.Join(userHomeDir, ".rk/logs")
	//os.MkdirAll(rkDir, os.ModePerm)
	//
	//bytes := []byte(fmt.Sprintf(confStr, path.Join(rkDir, "query.log")))

	bytes := []byte(fmt.Sprintf(confStr, "stdout"))
	logger, _, _ = rklogger.NewZapLoggerWithBytes(bytes, rklogger.JSON)
	factory = rkquery.NewEventFactory(
		rkquery.WithAppName("rk-cmd"),
		rkquery.WithZapLogger(logger))
}

func GetEvent(op string) rkquery.Event {
	event := factory.CreateEvent()
	event.SetOperation(op)
	event.SetStartTime(time.Now())
	return event
}

func GetEventV2(ctx *cli.Context) rkquery.Event {
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

func Success() {
	color.Green("✅ [success]")
	color.White("--------------------------------")
}

func Error(event rkquery.Event, err error) error {
	if err != nil {
		color.Red("[Error] %s", err.Error())
	} else {
		color.Red("[Error] internal error occur")
	}

	Finish(event, err)
	return err
}

func ErrorV2(err error) error {
	if err != nil {
		color.Red("[Error] %s", err.Error())
	} else {
		color.Red("[Error] internal error occur")
	}

	return err
}

// marshal to file to struct
func MarshalBuildConfig(path string, dest interface{}, event rkquery.Event) error {
	// 1: read file first
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		wrap := NewFileOperationError(err.Error())
		return Error(event, wrap)
	}

	if err := yaml.Unmarshal(bytes, dest); err != nil {
		wrap := NewUnMarshalError(err.Error())
		return Error(event, wrap)
	}

	event.AddPair("marshal-config", "success")
	return nil
}

// clear folder
func ClearTargetFolder(path string, event rkquery.Event) {
	// just try our best to remove it
	if err := os.RemoveAll(path); err != nil {
		color.Red(err.Error())
		event.AddPair("clear-target", "fail")
	}

	event.AddPair("clear-target", "success")
}

// validate go installation
func ValidateCommand(command string, event rkquery.Event) error {
	bytes, err := exec.Command("which", command).CombinedOutput()
	if err != nil {
		return Error(event,
			NewInvalidEnvError(
				fmt.Sprintf("%s is missing, please install go environment from offical site \n[err] %v\n[stderr] %s", command, err, string(bytes))))
	}

	color.Yellow(string(bytes))
	event.AddPair(fmt.Sprintf("%s-check", command), "success")

	return nil
}

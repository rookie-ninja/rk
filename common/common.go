// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
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
	factory *rk_query.EventFactory
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
	userHomeDir, _ := os.UserHomeDir()
	rkDir := path.Join(userHomeDir, ".rk/logs")
	os.MkdirAll(rkDir, os.ModePerm)

	bytes := []byte(fmt.Sprintf(confStr, path.Join(rkDir, "query.log")))

	logger, _, _ = rk_logger.NewZapLoggerWithBytes(bytes, rk_logger.JSON)
	factory = rk_query.NewEventFactory(
		rk_query.WithAppName("rk-cmd"),
		rk_query.WithLogger(logger))
}

func GetEvent(op string) rk_query.Event {
	event := factory.CreateEvent()
	event.SetOperation(op)
	event.SetStartTime(time.Now())
	return event
}

func Finish(event rk_query.Event, err error) {
	if err != nil {
		event.AddErr(err)
		event.SetResCode("Failed")
	} else {
		event.SetResCode("Success")
	}

	event.SetEventId(uuid.New().String())
	event.SetEndTime(time.Now())
	event.WriteLog()
}

func Success() {
	color.Green("[success]")
	color.White("--------------------------------")
}

func Error(event rk_query.Event, err error) error {
	color.Red(err.Error())
	Finish(event, err)
	return err
}

// marshal to file to struct
func MarshalBuildConfig(path string, dest interface{}, event rk_query.Event) error {
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
func ClearTargetFolder(path string, event rk_query.Event) {
	// just try our best to remove it
	if err := os.RemoveAll(path); err != nil {
		color.Red(err.Error())
		event.AddPair("clear-target", "fail")
	}

	event.AddPair("clear-target", "success")
}

// validate go installation
func ValidateCommand(command string, event rk_query.Event) error {
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

// Copyright (c) 2020 rookie-ninja
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rk_common

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rookie-ninja/rk-logger"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"os"
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

// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE
//
package xlog

import (
	"path"
	"runtime"
	"strconv"

	"github.com/krossdev/iam-ms/mss/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// common field names
const (
	FReqID  = "reqid"  // http request id
	FIP     = "ip"     // http client ip
	FPath   = "path"   // http request url path
	FMethod = "method" // http request method
)

var X = logrus.StandardLogger()

// Setup logger
func Setup(debug bool, config *config.Log) {
	X.SetReportCaller(true)

	// Logrus log to file kiam.log
	rotate_logger := &lumberjack.Logger{
		Filename:  path.Join(config.Path, "kiam-ms.log"),
		MaxSize:   20,
		Compress:  true,
		LocalTime: true,
	}
	// log json format
	X.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			function = path.Base(f.Function)
			file = path.Base(f.File) + ":" + strconv.Itoa(f.Line)
			return
		},
	})
	X.SetOutput(rotate_logger)

	if debug {
		X.AddHook(newTerminalHook())
		X.SetLevel(logrus.TraceLevel)
	} else {
		X.SetLevel(logrus.InfoLevel)
	}

	// add airbrake hook if enabled
	if config.Airbrake.Pid != 0 && len(config.Airbrake.Key) > 0 {
		X.AddHook(newAirbrakeHook(debug, config))
	}
}

// F is a shortcut of withFields
//
//   withFields usage:  xlog.withFields(logrus.Fields{"key", value, ...}).Info(...)
//   F usage: xlog.F("key", value, ...).Info(...)
func F(args ...interface{}) *logrus.Entry {
	if len(args)%2 != 0 {
		X.Panicf("Number of F(...) args must be even, current %d", len(args))
	}
	fields := logrus.Fields{}

	for i := 0; i < len(args); i += 2 {
		if s, ok := args[i].(string); !ok {
			X.Panicf("#%d arg %[2]v(%[2]T) must be string", i, args[i])
		} else {
			fields[s] = args[i+1]
		}
	}
	return X.WithFields(fields)
}

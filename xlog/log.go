// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package xlog

import (
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

var X = logrus.StandardLogger()

// Setup logger
func Setup(debug bool, output io.Writer) {
	X.SetReportCaller(true)

	if debug {
		X.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				function = path.Base(f.Function)
				file = path.Base(f.File) + ":" + strconv.Itoa(f.Line)
				return
			},
		})
		writer := io.MultiWriter(output, os.Stdout)
		X.SetOutput(writer)
		X.SetLevel(logrus.TraceLevel)
	} else {
		X.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				function = path.Base(f.Function)
				file = path.Base(f.File) + ":" + strconv.Itoa(f.Line)
				return
			},
		})
		X.SetOutput(output)
		X.SetLevel(logrus.InfoLevel)
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

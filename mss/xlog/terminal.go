// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE
//
package xlog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type terminalHook struct {
}

func newTerminalHook() *terminalHook {
	return &terminalHook{}
}

func (h *terminalHook) Fire(entry *logrus.Entry) error {
	colors := map[logrus.Level]color.Attribute{
		logrus.DebugLevel: color.FgHiBlack,
		logrus.TraceLevel: color.FgHiBlack,
		logrus.InfoLevel:  color.FgCyan,
		logrus.WarnLevel:  color.FgYellow,
		logrus.ErrorLevel: color.FgRed,
		logrus.FatalLevel: color.FgMagenta,
		logrus.PanicLevel: color.FgHiRed,
	}
	print := color.New(colors[entry.Level]).SprintfFunc()
	level := print("[%s]", entry.Level.String())

	var message string

	err, ok := entry.Data["error"].(error)
	if ok {
		message = fmt.Sprintf("%s %s: %v", level, entry.Message, err)
	} else {
		message = fmt.Sprintf("%s %s", level, entry.Message)
	}
	_, err = fmt.Println(message)
	return err
}

func (h *terminalHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

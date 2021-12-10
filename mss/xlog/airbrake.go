// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE
//
package xlog

import (
	"errors"
	"fmt"

	"github.com/krossdev/iam-ms/mss/config"

	"github.com/airbrake/gobrake/v5"
	"github.com/sirupsen/logrus"
)

type airbrakeHook struct {
	notifier *gobrake.Notifier
}

func newAirbrakeHook(debug bool, config *config.Log) *airbrakeHook {
	notifier := gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:   config.Airbrake.Pid,
		ProjectKey:  config.Airbrake.Key,
		Environment: "ms",
	})

	return &airbrakeHook{notifier: notifier}
}

func (h *airbrakeHook) Fire(entry *logrus.Entry) error {
	var noticeErr error

	err, ok := entry.Data["error"].(error)
	if ok {
		noticeErr = fmt.Errorf("%s: %v", entry.Message, err)
	} else {
		noticeErr = errors.New(entry.Message)
	}
	notice := h.notifier.Notice(noticeErr, nil, 6)

	for k, v := range entry.Data {
		notice.Context[k] = fmt.Sprintf("%s", v)
	}

	// ensure that logs are delivered before the process exits when panic
	if entry.Level == logrus.PanicLevel {
		h.notifier.SendNotice(notice)
	} else {
		h.notifier.SendNoticeAsync(notice)
	}
	return nil
}

func (h *airbrakeHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

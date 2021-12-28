package log

import (
	"github.com/krossdev/iam-ms/msc"
)

type LogRequest struct {
	hint string
	text string
}

func Publish(hint string, text string) error {
	b := msc.Borker()
	if b == nil {
		msc.Logger.Panicf("broker is nil")
	}
	rq := LogRequest{
		hint: hint,
		text: text,
	}
	return b.PublishLog(msc.SubjectLog, &rq)
}

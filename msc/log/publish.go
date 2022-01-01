package log

import (
	"github.com/krossdev/iam-ms/msc"
)

func Publish(rq interface{}) error {
	b := msc.Borker()
	if b == nil {
		msc.Logger.Panicf("broker is nil")
	}
	return b.Publish(msc.SubjectLog, rq)
}

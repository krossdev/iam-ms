// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"
	"time"

	"github.com/krossdev/iam-ms/msc"
)

type protoError struct {
	code int32
	err  error
}

func (e *protoError) Error() string {
	return e.err.Error()
}

func newProtoError(code int32, err error) *protoError {
	return &protoError{code, err}
}

func checkRequest(qt *msc.Request) *protoError {
	// check request version
	if err := checkVersion(qt.Version); err != nil {
		return newProtoError(msc.ReplyBadVersion, err)
	}
	// check request time
	td := time.Since(time.UnixMicro(qt.Time))
	if td < 0 || td > 10*time.Second {
		err := fmt.Errorf("request time %d invalid", qt.Time)
		return newProtoError(msc.ReplyBadTime, err)
	}
	// request must has reqid
	if len(qt.ReqId) == 0 {
		err := fmt.Errorf("request missing reqid")
		return newProtoError(msc.ReplyNoReqid, err)
	}
	return nil
}

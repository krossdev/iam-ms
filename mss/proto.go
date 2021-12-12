// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"
	"time"

	"github.com/blang/semver/v4"
	"github.com/krossdev/iam-ms/msc"
	"github.com/pkg/errors"
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

// check request version is compatible
func checkVersion(version string) error {
	if len(version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}
	rv, err := semver.Parse(version) // request version
	if err != nil {
		return errors.Wrap(err, "parse request version error")
	}
	lv := semver.MustParse(msc.Version) // server version

	// major version must equal
	if rv.Major != lv.Major {
		return fmt.Errorf(
			"version incompatible, expect %v, got %v", msc.Version, version,
		)
	}
	// minor version must equal if major version is 0
	if rv.Major == 0 && rv.Minor != lv.Minor {
		return fmt.Errorf(
			"version incompatible, expect %v, got %v", msc.Version, version,
		)
	}
	return nil
}

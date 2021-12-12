// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// request data
type Request struct {
	Version string      `json:"version"` // protocol version
	Time    int64       `json:"time"`    // request timestamp, in microsecond
	ReqId   string      `json:"reqid"`   // request id, unique
	Action  string      `json:"action"`  // action
	Payload interface{} `json:"payload"` // payload
}

// reply data
type Reply struct {
	Version string      `json:"version"` // protocol version
	Time    int64       `json:"time"`    // response timestamp, in microsecond
	ReqId   string      `json:"reqid"`   // request id, unique
	Code    int32       `json:"code"`    // response code
	Message string      `json:"message"` // response message
	Payload interface{} `json:"payload"` // payload
}

// reply code
const (
	ReplyOk         = 0
	ReplyNoReply    = 100
	ReplyBadVersion = 101
	ReplyBadTime    = 102
)

// make a request package
func makeRequest(action string, payload interface{}) *Request {
	return &Request{
		Version: Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   strings.ToLower(strings.ReplaceAll(uuid.NewString(), "-", "")),
		Action:  action,
		Payload: payload,
	}
}

// check reply
func checkReply(rp *Reply, reqid string) error {
	if len(rp.Version) == 0 {
		return fmt.Errorf("reply missing version")
	}
	if len(rp.ReqId) == 0 {
		return fmt.Errorf("reply missing reqid")
	}
	// reply reqid must match to request reqid
	if rp.ReqId != reqid {
		return fmt.Errorf("reply reqid '%s' dismatch to '%s'", rp.ReqId, reqid)
	}
	return nil
}

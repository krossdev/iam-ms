// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// data send to server
type RequestAction struct {
	Version   string      `json:"version"`   // protocol version
	Timestamp int64       `json:"timestamp"` // request timestamp, in micro-second
	ReqId     string      `json:"reqid"`     // request id, unique
	Subject   string      `json:"subject"`   // subscribe subject
	Name      string      `json:"name"`      // action name
	Payload   interface{} `json:"payload"`   // payload
}

// data response from server
type ReplyAction struct {
	Version   string      `json:"version"`   // protocol version
	Timestamp int64       `json:"timestamp"` // response timestamp, in micro-second
	ReqId     string      `json:"reqid"`     // request id, unique
	Code      int32       `json:"code"`      // response code
	Message   string      `json:"message"`   // response message
	Payload   interface{} `json:"payload"`   // payload
}

// server response code
const (
	ReplyCodeOk                  = 0
	ReplyCodeVersionIncompatible = 101
)

// send a request to server and receive a reply
func sendActionRequest(a *RequestAction) (interface{}, error) {
	var reply interface{}

	a.Version = Version
	a.Timestamp = time.Now().UnixMicro()
	a.ReqId = strings.ToLower(strings.ReplaceAll(uuid.NewString(), "-", ""))

	err := conn.Request(a.Subject, a, &reply, 10*time.Second)
	if err != nil {
		return nil, err
	}
	logger.Infof("reply: %+v\n", reply)
	return &reply, nil
}

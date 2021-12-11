// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package msc

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// data send to server
type RequestAction struct {
	Version string      `json:"version"` // protocol version
	Time    int64       `json:"time"`    // request timestamp, in micro-second
	ReqId   string      `json:"reqid"`   // request id, unique
	Subject string      `json:"subject"` // subscribe subject
	Name    string      `json:"name"`    // action name
	Payload interface{} `json:"payload"` // payload
}

// data response from server
type ReplyAction struct {
	Version string      `json:"version"` // protocol version
	Time    int64       `json:"time"`    // response timestamp, in micro-second
	ReqId   string      `json:"reqid"`   // request id, unique
	Code    int32       `json:"code"`    // response code
	Message string      `json:"message"` // response message
	Payload interface{} `json:"payload"` // payload
}

// server response code
const (
	ReplyCodeOk         = 0
	ReplyCodeNoReply    = 100
	ReplyCodeBadVersion = 101
	ReplyCodeBadTime    = 102
)

// send a request to server and receive a reply
func sendActionRequest(q *RequestAction) (interface{}, error) {
	var reply ReplyAction

	if len(q.Name) == 0 {
		logger.Panicf("request action missing name")
	}
	if len(q.Subject) == 0 {
		logger.Panicf("request action missing subject")
	}
	q.Version = Version
	q.Time = time.Now().UnixMicro()
	q.ReqId = strings.ToLower(strings.ReplaceAll(uuid.NewString(), "-", ""))

	l := logger.WithFields(logrus.Fields{
		"version": q.Version,
		"time":    q.Time,
		"reqid":   q.ReqId,
		"name":    q.Name,
	})
	l.Tracef("action %s send to %s", q.Name, q.Subject)

	// send the request to broker
	err := conn.Request(q.Subject, q, &reply, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "request error")
	}
	if err = checkReplyAction(&reply); err != nil {
		return nil, errors.Wrap(err, "reply invalid")
	}
	// reqid must match
	if reply.ReqId != q.ReqId {
		err = fmt.Errorf("reply reqid '%s' dismatch to '%s'", reply.ReqId, q.ReqId)
		return nil, err
	}
	dt := reply.Time - q.Time

	l.WithField("code", reply.Code).Tracef(
		"action %s reply with code %d, latency %fms",
		q.Name, reply.Code, float64(dt/1000),
	)

	if reply.Code != ReplyCodeOk {
		return nil, fmt.Errorf("%s", reply.Message)
	}
	return &reply, nil
}

// check reply fields
func checkReplyAction(r *ReplyAction) error {
	if len(r.Version) == 0 {
		return fmt.Errorf("missing version")
	}
	if len(r.ReqId) == 0 {
		return fmt.Errorf("missing reqid")
	}
	return nil
}

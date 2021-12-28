package msc

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// semantic versioning, https://semver.org
const Version = "0.1.1"

const (
	SubjectAction = "kiam.console.action"
	SubjectNotice = "kiam.console.notice"
	SubjectAudit  = "kiam.console.audit"
	SubjectLog    = "kiam.console.log"
)

// request data
type Request struct {
	Version string      `json:"version"` // protocol version
	Time    int64       `json:"time"`    // request timestamp, in microsecond
	ReqId   string      `json:"reqid"`   // request id, unique
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
	ReplyOk         = 0   // ok
	ReplyError      = 1   // general error
	ReplyNoReply    = 100 // no reply(inbox) subject present
	ReplyBadVersion = 101 // version incompatibled
	ReplyBadTime    = 102 // time incorrent
	ReplyNoReqid    = 103 // missing reqid
	ReplyNoAction   = 104 // missing action
	ReplyNotImp     = 105 // not implemented
)

// make a request package
func makeRequest(payload interface{}) *Request {
	return &Request{
		Version: Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   strings.ToLower(strings.ReplaceAll(uuid.NewString(), "-", "")),
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

// make a reply package
func MakeReply(code int32, message string, payload interface{}, reqid string) *Reply {
	return &Reply{
		Version: Version,
		Time:    time.Now().UnixMicro(),
		ReqId:   reqid,
		Code:    code,
		Message: message,
		Payload: payload,
	}
}

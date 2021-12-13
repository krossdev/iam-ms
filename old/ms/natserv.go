// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package ms

import (
	"encoding/json"
	"github.com/krossdev/iam-ms/config"
	"github.com/krossdev/iam-ms/email"
	"github.com/krossdev/iam-ms/util"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"log"
	"time"
)

var ec *nats.EncodedConn
var nc *nats.Conn

func connect() (err error) {
	c := config.G.Nats
	nc, err = nats.Connect(c.Url,
		nats.Name("iam-ms"),
		nats.Timeout(10*time.Second),
		nats.ReconnectWait(1*time.Second),  // wait 1 second before reconnect
		nats.ReconnectBufSize(8*1024*1024), // cache during reconnecting
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("nats disconnect: %s", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("nats reconnected")
		}),
	)
	if err != nil {
		return
	}
	err = nc.Flush()
	if err != nil {
		return
	}
	// err = nc.LastError()
	ec, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	return
}

func mailMsgHandler(m *email.Mail) error {
	err := email.Send(m)
	return err
}

func mailQueueSubscribe2() {
	cb := func(msg *nats.Msg) {
		// in nats goroutine, and be always the same goroutine
		var m *email.Mail
		err := json.Unmarshal(msg.Data, &m)
		if err != nil {
			err = errors.WithStack(err)
		} else {
			err = mailMsgHandler(m)
		}
		if err != nil {
			log.Printf("mail send error: %+v", errors.WithStack(err))
		}
		if m.Reply {
			// cli need reply
			if err != nil {
				err = msg.Respond([]byte(err.Error()))
			} else {
				var ok []byte
				err = msg.Respond(ok)
			}
			if err != nil {
				log.Printf("nats reply error: %+v", errors.WithStack(err))
			}
		}
	}
	if _, err := nc.QueueSubscribe("mail", "workers", cb); err != nil {
		err = errors.WithStack(err)
		log.Fatalf("QueueSubscribe error: %+v", err)
	}
}

func mailQueueSubscribe() *nats.Subscription {
	handler := func(m *email.Mail) error {
		err := email.Send(m)
		if err != nil {
			log.Printf("mail send error: %+v", err)
		}
		return err
	}
	sub, err := ec.QueueSubscribe("mail", "workers", func(m *email.Mail) {
		// in nats goroutine, always the same one
		log.Printf("sub-callback goid = %d", util.CurGoroutineID())
		//go func() { // 已经位于 nats subscribe routine, 没有必要再启动 routine
		log.Printf("sub-dispath goid = %d, %v", util.CurGoroutineID(), m)
		err := handler(m)
		if err != nil {
		}
		//}()
	})
	if err != nil {
		log.Fatal(err)
	}
	return sub
}

func Setup() (err error) {
	err = connect()
	if err != nil {
		return errors.Wrap(err, "connect to nats error")
	} else {
		log.Printf("nats connected")
	}

	mailQueueSubscribe2()
	log.Printf("sub-setup goid = %d", util.CurGoroutineID())

	return
}

// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package mscli

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Params struct {
	Url     string
	Timeout time.Duration
}

var p *Params

func Setup(params *Params) error {
	p = params

	// set default params
	if p.Timeout == 0 {
		p.Timeout = 10
	}

	return setup(p.Url)
}

var ec *nats.EncodedConn
var nc *nats.Conn

func connect(natsUrl string) (err error) {
	nc, err = nats.Connect(natsUrl,
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
	ec, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	return
}

func sendMail(m *Mail) (err error) {
	var msg []byte
	if msg, err = json.Marshal(m); err != nil {
		return
	}
	if m.Reply {
		var reply *nats.Msg
		reply, err = nc.Request("mail", msg, p.Timeout*time.Second)
		if err == nil {
			if len(reply.Data) > 0 {
				// reply error
				err = errors.New(string(reply.Data))
			}
		}
	} else {
		err = nc.Publish("mail", msg)
	}
	if m.AfterHook != nil {
		err = m.AfterHook(m, err)
	}
	return
}

func setup(natsUrl string) (err error) {
	if err = connect(natsUrl); err != nil {
		log.Printf("connect to nats error: %s", err)
		return err
	}
	// defer ec.Conn.Close()

	go func() {
		// loop infinite
		for {
			select {
			case m := <-mailC:
				err = sendMail(m)
				if err != nil {
					log.Printf("goroutine select loop: send mail to nats error: %s", err)
				}
			}
		}
	}()

	return
}

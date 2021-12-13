// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/email"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/pkg/errors"
)

// command line options
var (
	configFile = flag.String("c", "./config.yaml", "Configuration `filepath`")
)

func main() {
	// parse command line options
	flag.Parse()

	// setup everythings
	if err := load(); err != nil {
		xlog.X.WithError(err).Fatalf("Startup failure")
	}

	// wait for shutdown
	wg := sync.WaitGroup{}
	wg.Add(1)

	// catch signal to reload or graceful showdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for s := range c {
			if s == syscall.SIGHUP {
				xlog.X.Info("receive SIGHUP, try to reload...")

				if err := load(); err != nil {
					xlog.X.WithError(err).Error("reload failed")
				}
				continue
			}
			xlog.X.Infof("receive %s signal, shutdown now...", s.String())

			disconnect() // disconnect from message broker
			wg.Done()
			break
		}
	}()

	xlog.X.Infof("%d is ready, wait for incoming message...", os.Getpid())
	wg.Wait()
}

// (re)load server
// when receive SIGHUP will call this function to reload server
func load() error {
	// parse configuration file
	conf, err := config.Load(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}
	// setup log
	xlog.Setup(conf.Debug, &conf.Log)

	// set msc logger
	msc.SetLogger(xlog.F(xlog.FRealm, "ms"))

	// propagation mail configuration
	email.Setup(&conf.Mail)

	// connect to message border
	if err = reconnect(conf.Brokers); err != nil {
		return errors.Wrap(err, "failed to connect message borker")
	}

	// subscribe action subject
	if err = subscribeAction(); err != nil {
		return errors.Wrap(err, "failed to subscribe action subject")
	}
	return nil
}

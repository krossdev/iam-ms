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
	"github.com/krossdev/iam-ms/mss/xlog"
	"github.com/pkg/errors"
)

// command line options
var (
	configFile = flag.String("c", "./config.yaml", "Configuration `filepath`")
	reload     = flag.Bool("r", false, "Reload server")
)

func main() {
	// parse command line options
	flag.Parse()

	// setup everythings
	if err := load(); err != nil {
		xlog.X.Fatalf("Startup failure")
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

	xlog.X.Info("Ready, wait for incoming message...")
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
	msc.SetLogger(xlog.F(xlog.FRealm, "ms"))

	// connect to message border
	if err = reconnect(conf.Brokers); err != nil {
		return errors.Wrap(err, "failed to connect message borker")
	}

	if err = subscribeConsoleActions(); err != nil {
		return errors.Wrap(err, "failed to subscribe console actions")
	}
	return nil
}

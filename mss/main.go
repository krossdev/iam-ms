// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/krossdev/iam-ms/msc"
	"github.com/krossdev/iam-ms/mss/config"
	"github.com/krossdev/iam-ms/mss/xlog"
)

// command line options
var (
	configFile = flag.String("c", "./config.yaml", "Configuration `filepath`")
)

func main() {
	// parse command line options
	flag.Parse()

	// parse configuration file
	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configration: %s", err)
	}

	// setup xlog
	xlog.Setup(conf.Debug, &conf.Log)

	// msc logger
	msc.SetLogger(xlog.F(xlog.FRealm, "ms"))

	// connect to message border
	if err = connect(conf.Brokers); err != nil {
		xlog.X.Fatalf("Failed to connect message broker: %v", err)
	}

	if err = subscribeConsoleActions(); err != nil {
		log.Fatalf("Failed to subscribe console actions")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Catch signal to graceful showdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Stop the server when signal arrived(exclude SIGHUP)
	go func() {
		s := <-quit
		if s == syscall.SIGHUP {
			xlog.X.Info("receive SIGHUP, ignore...")
			return
		}
		xlog.X.Infof("receive %s signal, shutdown now...", s.String())

		// disconnect from message broker
		disconnect()
		wg.Done()
	}()

	xlog.X.Info("Ready, wait for incoming message...")
	wg.Wait()
}

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
	"github.com/krossdev/iam-ms/mss/geoip"
	"github.com/krossdev/iam-ms/mss/xlog"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

// command line options
var (
	configFile = flag.String("c", "./config.yaml", "Configuration `filepath`")
	watch      = flag.Bool("watch", false, "Watch configuration file change to reload")
)

func main() {
	// parse command line options
	flag.Parse()

	// setup everythings
	if err := load(); err != nil {
		xlog.X.WithError(err).Fatalf("Startup failure")
	}

	// watch configuration file change to reload
	if *watch {
		go watchConfig(*configFile)
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

				// reload server
				if err := load(); err != nil {
					xlog.X.WithError(err).Error("reload failed")
				}
				continue
			}
			xlog.X.Infof("receive %s signal, shutdown now...", s.String())

			disconnect() // disconnect from message broker
			wg.Done()    // quit the app
			break
		}
	}()

	xlog.X.Infof("%d is ready, wait for incoming messages...", os.Getpid())
	wg.Wait()
}

// (re)load server
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

	// propagation geoip configuration
	geoip.Setup(&conf.Geoip)

	// connect to message border
	if err = reconnect(conf.Brokers); err != nil {
		return errors.Wrap(err, "failed to connect message borker")
	}

	// subscribe actions
	if err = subscribeActionsWithConfig(&conf.Service.Actions); err != nil {
		return errors.Wrap(err, "failed to subscribe action subject")
	}
	// subscribe audit
	if err = subscribeAuditsWithConfig(&conf.Service.Audits); err != nil {
		return errors.Wrap(err, "failed to subscribe audit subject")
	}
	return nil
}

// watch config file change to reload
func watchConfig(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		xlog.X.WithError(err).Error("watch config failed")
	}
	defer watcher.Close()

	if err = watcher.Add(file); err != nil {
		xlog.X.WithError(err).Error("watch config failed")
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				xlog.X.Info("configuration file has changed, try to reload...")

				if err := load(); err != nil {
					xlog.X.WithError(err).Error("reload failed")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			xlog.X.WithError(err).Error("watch config error")
		}
	}
}

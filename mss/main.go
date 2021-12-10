// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"flag"
	"log"

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

	// connect to message border
	if err = connect(conf.Brokers); err != nil {
		xlog.X.Fatalf("Failed to connect message broker: %v", err)
	}
}

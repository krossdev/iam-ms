// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package main

import (
	"flag"
	"github.com/krossdev/iam-ms/config"
	"github.com/krossdev/iam-ms/ms"
	"log"
)

func _t(key string, args ...interface{}) string {
	return key
}

// Command line options
var (
	configFile = flag.String("c", "./conf/config.yaml", _t("Configuration `filepath`"))
	//secureFile = flag.String("s", "", _t("Secure configuration `filepath`"))
	//addUser    = flag.Bool("adduser", false, _t("Add the first console user"))
)

func main() {
	// Parse command line options
	flag.Parse()

	// Parse configuration file
	if err := config.LoadConfig(*configFile); err != nil {
		log.Fatalf("Unable to load configration: %+v", err)
	}

	// start ms server
	err := ms.Setup()
	if err != nil {
		log.Fatalf("ms setup error %+v", err)
	} else {
		log.Println("ms server started.")
	}

	// wait for quit?
	for {
		select {}
	}
}

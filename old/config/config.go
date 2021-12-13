// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package config

import (
	"github.com/krossdev/iam-ms/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// G Global config instance
var G *Configuration

type Configuration struct {
	RootDir string // conf root directory
	Nats    struct {
		Url string
	}
	Mail struct {
		SchemaDir string
	}
}

// LoadConfig load configuration file
func LoadConfig(filename string) error {
	if G != nil {
		// return fast
		return nil
	}
	var file []byte
	var err error
	file, err = os.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "load conf file `%s`", filename)
	}
	if err = yaml.Unmarshal(file, &G); err != nil {
		return errors.Wrapf(err, "when parse conf file `%s`", filename)
	}

	G.RootDir = filepath.Dir(filename)

	// check email conf
	if G.Mail.SchemaDir == "" {
		G.Mail.SchemaDir = G.RootDir + "/schema"
	}

	// nats conf
	if G.Nats.Url == "" {
		log.Fatal("nats url not set")
	}

	return nil
}

// LoadTestConfig just for test
func LoadTestConfig() {
	if G == nil {
		err := LoadConfig(getTestConfigFile())
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// just for test
func getTestConfigFile() string {
	_, file, _, _ := runtime.Caller(0)
	for d := file; d != "" && d != "/"; d = filepath.Dir(d) {
		filename := filepath.Join(d, "conf/config.yaml")
		if util.MustFileExists(filename) {
			return filename
		}
	}
	log.Fatalf("%+v", errors.New("cannot find test conf"))
	return ""
}

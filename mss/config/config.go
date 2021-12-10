// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE
//
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Log struct {
	Path     string `yaml:"path"`
	Airbrake struct {
		Pid int64  `yaml:"pid"`
		Key string `yaml:"key"`
	}
}

type Configuration struct {
	Debug   bool     `yaml:"debug"`
	Log     Log      `yaml:"log"`
	Brokers []string `yaml:"brokers"`
}

// Load configuration from file
func Load(path string) (*Configuration, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Configuration

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	// check required configuration entries
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("brokers cannot be empty")
	}
	return &config, nil
}

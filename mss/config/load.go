// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE
//
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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
	if len(config.Log.Path) == 0 {
		return nil, fmt.Errorf("log.path cannot be empty")
	}
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("brokers cannot be empty")
	}
	return &config, nil
}

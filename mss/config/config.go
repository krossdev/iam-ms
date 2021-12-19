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

type ActionIPLocation struct {
	Enabled bool   `yaml:"enabled"`
	Engine  string `yaml:"engine"`
}

const (
	IPLocationEngineGeoip = "geoip"
)

type ActionSendVerifyEmail struct {
	Enabled bool `yaml:"enabled"`
}

type ServiceActions struct {
	IPLocation      ActionIPLocation      `yaml:"ip_location"`
	SendVerifyEmail ActionSendVerifyEmail `yaml:"send_verify_email"`
}

// Service
type Service struct {
	Event   bool           `yaml:"event"`
	Audit   bool           `yaml:"audit"`
	Actions ServiceActions `yaml:"actions"`
}

type Mta struct {
	Name    string   `yaml:"name"`
	Host    string   `yaml:"host"`
	Port    int      `yaml:"port"`
	SSLMode bool     `yaml:"sslmode"`
	Sender  string   `yaml:"sender"`
	ReplyTo string   `yaml:"replyto"`
	CC      []string `yaml:"cc"`
	BCC     []string `yaml:"bcc"`
	User    string   `yaml:"user"`
	Passwd  string   `yaml:"passwd"`
}

// Mail
type Mail struct {
	SubjectPrefix string `yaml:"subject_prefix"`
	TemplateDir   string `yaml:"template_dir"`
	Mtas          []Mta  `yaml:"mtas"`
	PreferredMta  string `yaml:"preferred_mta"`
}

type Geoip struct {
	Path string `yaml:"path"`
}

type Configuration struct {
	Debug   bool     `yaml:"debug"`
	Log     Log      `yaml:"log"`
	Brokers []string `yaml:"brokers"`
	Service Service  `yaml:"service"`
	Mail    Mail     `yaml:"mail"`
	Geoip   Geoip    `yaml:"geoip"`
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
	if len(config.Log.Path) == 0 {
		return nil, fmt.Errorf("log.path cannot be empty")
	}
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("brokers cannot be empty")
	}
	return &config, nil
}

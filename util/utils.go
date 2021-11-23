// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package util

import (
	"os"
)

// E2P error to panic
func E2P(err error) {
	if err != nil {
		panic(err)
	}
}

// MustGetwd wrap os.Getwd
func MustGetwd() string {
	wd, err := os.Getwd()
	E2P(err)
	return wd
}

// MustExists check path exist
func MustExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

// MustFileExists check path exist and type file
func MustFileExists(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return !info.IsDir()
	}
	return false
}

// MustDirExists check path exist and type dir
func MustDirExists(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return info.IsDir()
	}
	return false
}

// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package main

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/pkg/errors"

	"github.com/krossdev/iam-ms/msc"
)

// check request version is compatible
func checkVersion(version string) error {
	if len(version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}
	rv, err := semver.Parse(version) // request version
	if err != nil {
		return errors.Wrap(err, "parse request version error")
	}
	lv := semver.MustParse(msc.Version) // server version

	// major version must equal
	if rv.Major != lv.Major {
		return fmt.Errorf(
			"version incompatible, expect %v, got %v", msc.Version, version,
		)
	}
	// minor version must equal if major version is 0
	if rv.Major == 0 && rv.Minor != lv.Minor {
		return fmt.Errorf(
			"version incompatible, expect %v, got %v", msc.Version, version,
		)
	}
	return nil
}

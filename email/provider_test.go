// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package email

import (
	"errors"
	"fmt"
	"testing"
)

func xx1() (ps []*provider, err error) {
	p := provider{Name: "xx"}
	fmt.Println("ps is nil?", ps == nil)

	fmt.Println("init ps=", ps)
	return append(ps, &p, &p), errors.New("hello")
}
func Test1(t *testing.T) {
	a, b := xx1()
	fmt.Println(a, "-", b)
}

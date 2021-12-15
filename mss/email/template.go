// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package email

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type TemplateData struct {
	Logo  template.URL // logo url, can be a http url, a cid:xxx or a base64 dataurl
	Title string       // title
}

const (
	TemplateVerifyEmail = "verify-email"
)

// load template from file, support i18n
func ExecTemplate(tname string, locale string, data interface{}) (string, error) {
	var tpath string

	// lookup locale template first
	if len(locale) > 0 {
		lname := fmt.Sprintf("%s_%s.html", tname, strings.ToLower(locale))
		lpath := path.Join(mailConfig.TemplateDir, lname)
		if _, err := os.Stat(lpath); err == nil {
			tpath = lpath
		}
	}
	// no locale template found, fallback to default template
	if len(tpath) == 0 {
		tpath = path.Join(mailConfig.TemplateDir, fmt.Sprintf("%s.html", tname))
	}

	t, err := template.ParseFiles(tpath)
	if err != nil {
		return "", errors.Wrap(err, "parse template error")
	}
	var out bytes.Buffer

	if err = t.Execute(&out, data); err != nil {
		return "", errors.Wrap(err, "execute template error")
	}
	return out.String(), nil
}

// return image path in template dir
func ImagePath(name string) string {
	return path.Join(mailConfig.TemplateDir, name)
}

// return logo path in template dir
func LogoPath() string {
	return ImagePath("logo.png")
}

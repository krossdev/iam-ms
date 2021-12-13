// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam-ms/blob/main/LICENSE
//
package email

import (
	"bytes"
	"html/template"
	"path"
	"strings"

	"github.com/pkg/errors"
)

type TemplateData struct {
	Logo  template.URL // logo url, can be a http url, a cid:xxx or a base64 dataurl
	Title string       // title
}

// load template from file
func ExecTemplate(tname string, data interface{}) (string, error) {
	if !strings.HasSuffix(tname, ".html") {
		tname += ".html"
	}
	tpath := path.Join(mailConfig.TemplateDir, tname)

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

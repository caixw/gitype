// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTemplate(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{Theme: "t1"}}
	a.NotError(d.loadTemplate("./testdata/themes")).NotNil(d.Template)

	w := bytes.NewBufferString("")
	a.NotError(d.Template.ExecuteTemplate(w, "header.html", map[string]string{"val": "123"}))
	a.Equal(w.String(), "<h1>123</h1>\n")
}

func TestGetThemesName(t *testing.T) {
	a := assert.New(t)

	themes, err := getThemesName("./testdata/themes")
	a.NotError(err).NotNil(themes).Equal(2, len(themes))
}

func TestData_longDataFormat(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{LongDateFormat: "2006-01-02T15:04:05-0700"}}
	a.Equal(d.longDateFormat(1456324895), "2016-02-24T22:41:35+0800")
}

func TestData_shortDataFormat(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{ShortDateFormat: "2006-01-02"}}
	a.Equal(d.shortDateFormat(1456324895), "2016-02-24")
}

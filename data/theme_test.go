// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/gitype/data/loader"
	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	conf := &loader.Config{Theme: "t1"}
	theme, err := loadTheme(testdataPath, conf)
	a.NotError(err).NotNil(theme)

	a.Equal(theme.Name, "name")
	a.Equal(theme.Author.Name, "caixw")
}

func TestStripTags(t *testing.T) {
	a := assert.New(t)

	tests := map[string]string{
		"<div>str</div>":        "str",
		"str<br />":             "str",
		"<div><p>str</p></div>": "str",
	}

	for expr, val := range tests {
		a.Equal(stripTags(expr), val, "测试[%v]时出错", expr)
	}
}

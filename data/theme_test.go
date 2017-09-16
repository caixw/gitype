// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestLoadThemes(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("../testdata")

	ts, err := loadThemes(p)
	a.NotError(err).NotNil(ts).Equal(len(ts), 2)

	// 排序是否正常
	a.Equal(ts[0].ID, "t1")
	a.Equal(ts[1].ID, "t2")
}

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("../testdata")

	theme, err := loadTheme(p, "t1")
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

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"testing"

	"github.com/issue9/assert"
)

func TestIsIgnoreThemeFile(t *testing.T) {
	a := assert.New(t)
	a.True(isIgnoreThemeFile(templateExtension))
	a.True(isIgnoreThemeFile(".yaml"))
	a.False(isIgnoreThemeFile(".txt"))
	a.False(isIgnoreThemeFile(".css"))
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

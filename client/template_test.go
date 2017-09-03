// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestIsIgnoreThemeFile(t *testing.T) {
	a := assert.New(t)
	a.True(isIgnoreThemeFile(vars.TemplateExtension))
	a.True(isIgnoreThemeFile(".yaml"))
	a.False(isIgnoreThemeFile(".txt"))
	a.False(isIgnoreThemeFile(".css"))
}

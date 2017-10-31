// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package path

import (
	"testing"

	"github.com/caixw/gitype/vars"
	"github.com/issue9/assert"
)

func TestPath(t *testing.T) {
	a := assert.New(t)

	p := New("/")
	a.Equal(p.ConfDir, "/"+vars.ConfFolderName)

	// ThemesPath
	a.Equal(p.ThemesPath("def", "//style", "style.png"), "/data/themes/def/style/style.png")
	a.Equal(p.ThemesPath("def", "//style//style.png"), "/data/themes/def/style/style.png")
	a.Equal(p.ThemesPath("def", "//style//*.html"), "/data/themes/def/style/*.html")
}

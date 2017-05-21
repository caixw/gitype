// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := loadTheme("./testdata/themes", "t1")
	a.NotError(err).NotNil(theme)

	a.Equal(theme.Name, "name")
	a.Equal(theme.Author.Name, "name")
}

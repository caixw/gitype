// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/gitype/path"
	"github.com/issue9/assert"
)

var testdataPath = path.New("../testdata")

func TestLoad(t *testing.T) {
	a := assert.New(t)
	d, err := Load(testdataPath)
	a.NotError(err).NotNil(d)

	a.Equal(len(d.Posts), 2)

	// theme
	a.Equal(len(d.Themes), 2)
	a.Equal(d.Themes[0].ID, "t1") // 默认主题
	a.Equal(d.Themes[0].Author.Name, "caixw")

	// feed
	a.Equal(d.Opensearch.URL, "/opensearch.xml")
	a.Equal(d.Atom.URL, "/atom.xml")
	a.Nil(d.Sitemap)
}

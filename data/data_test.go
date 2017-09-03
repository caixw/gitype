// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestData_loadFiles(t *testing.T) {
	a := assert.New(t)

	d := &Data{path: vars.NewPath("./testdata")}
	a.NotError(d.loadFiles())
	a.NotNil(d.Tags).
		NotNil(d.Posts).
		NotNil(d.Themes).
		NotNil(d.Links)

	// Data.Tags
	a.Equal(d.Tags[0].Slug, "default1")
	a.Equal(d.Tags[0].Color, "efefef")
	a.Equal(d.Tags[0].Title, "默认1")
	a.Equal(d.Tags[1].Slug, "default2")
	a.Equal(d.Tags[0].Permalink, "") // 未调用 sanitize 初始化

	// Data.Links
	a.True(len(d.Links) > 0)
	a.Equal(d.Links[0].Text, "text0")
	a.Equal(d.Links[0].URL, "url0")
	a.Equal(d.Links[1].Text, "text1")
	a.Equal(d.Links[1].URL, "url1")
}

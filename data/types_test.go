// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestData_loadTags(t *testing.T) {
	a := assert.New(t)

	data := &Data{
		path: vars.NewPath("./testdata"),
	}
	a.NotError(data.loadTags())
	a.NotNil(data.Tags)
	a.Equal(data.Tags[0].Slug, "default1")
	a.Equal(data.Tags[0].Color, "efefef")
	a.Equal(data.Tags[0].Title, "默认1")
	a.Equal(data.Tags[1].Slug, "default2")
	a.Equal(data.Tags[0].Permalink, "/tags/default1.html")

	t.Log(data.Tags[0])
}

func TestData_loadLinks(t *testing.T) {
	a := assert.New(t)

	data := &Data{path: vars.NewPath("./testdata")}
	a.NotError(data.loadLinks())
	a.True(len(data.Links) > 0)
	a.Equal(data.Links[0].Text, "text0")
	a.Equal(data.Links[0].URL, "url0")
	a.Equal(data.Links[1].Text, "text1")
	a.Equal(data.Links[1].URL, "url1")

	t.Log(data.Links)
}

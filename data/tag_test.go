// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestData_loadTags(t *testing.T) {
	a := assert.New(t)

	data := &Data{URLS: &URLS{Root: "/root", Tag: "tags", Suffix: ".html"}}
	a.NotError(data.loadTags("./testdata/meta/tags.yaml"))
	a.NotNil(data.Tags)
	a.Equal(data.Tags[0].Slug, "default1")
	a.Equal(data.Tags[0].Color, "efefef")
	a.Equal(data.Tags[0].Title, "默认1")
	a.Equal(data.Tags[1].Slug, "default2")
	a.Equal(data.Tags[0].Permalink, "/root/tags/default1.html")

	t.Log(data.Tags[0])
}

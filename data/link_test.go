// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestData_loadLinks(t *testing.T) {
	a := assert.New(t)

	data := &Data{}
	a.NotError(data.loadLinks("./testdata/meta/links.yaml"))
	a.True(len(data.Links) > 0)
	a.Equal(data.Links[0].Text, "text0")
	a.Equal(data.Links[0].URL, "url0")
	a.Equal(data.Links[1].Text, "text1")
	a.Equal(data.Links[1].URL, "url1")

	t.Log(data.Links)
}

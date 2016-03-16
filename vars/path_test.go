// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"runtime"
	"strings"
	"testing"

	"github.com/issue9/assert"
)

func TestNew(t *testing.T) {
	// 仅检测unix-like系统，验证路径是否正确。
	// 由filepath保证其它系统上也能得到相同的结果。
	if runtime.GOOS == "windows" {
		return
	}

	a := assert.New(t)

	p, err := NewPath("./testdata")
	a.NotError(err).NotNil(p)

	eq := func(p1, p2 string) bool {
		return strings.HasSuffix(p1, p2)
	}

	a.True(eq(p.Root, "testdata")).
		True(eq(p.Data, "testdata/data")).
		True(eq(p.Conf, "testdata/conf")).
		True(eq(p.ConfApp, "testdata/conf/app.json")).
		True(eq(p.ConfLogs, "testdata/conf/logs.xml")).
		True(eq(p.DataConf, "testdata/data/meta/config.yaml")).
		True(eq(p.DataTags, "testdata/data/meta/tags.yaml")).
		True(eq(p.DataURLS, "testdata/data/meta/urls.yaml")).
		True(eq(p.DataPosts, "testdata/data/posts")).
		True(eq(p.DataThemes, "testdata/data/themes")).
		True(eq(p.DataRaws, "testdata/data/raws"))
}

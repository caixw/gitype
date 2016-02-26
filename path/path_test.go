// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package path

import (
	"runtime"
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

	p := New("./testdata")

	a.Equal(p.Root, "./testdata").
		Equal(p.Data, "testdata/data").
		Equal(p.Conf, "testdata/conf").
		Equal(p.ConfApp, "testdata/conf/app.json").
		Equal(p.ConfLogs, "testdata/conf/logs.xml").
		Equal(p.DataConf, "testdata/data/meta/config.yaml").
		Equal(p.DataTags, "testdata/data/meta/tags.yaml").
		Equal(p.DataURLS, "testdata/data/meta/urls.yaml").
		Equal(p.DataPosts, "testdata/data/posts").
		Equal(p.DataThemes, "testdata/data/themes").
		Equal(p.DataRaws, "testdata/data/raws")
}

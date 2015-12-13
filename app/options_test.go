// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_fromMaps(t *testing.T) {
	a := assert.New(t)

	maps := []map[string]string{
		{
			"key":   "pageSize",
			"group": "reading",
			"value": "50",
		},
		{ // 不存在的数据，会忽略
			"key":   "noexists",
			"group": "system",
			"value": "50",
		},
	}

	opt := &Options{}
	a.NotError(opt.fromMaps(maps))
	a.Equal(opt.PageSize, 50)
}

func TestOptions_setValue(t *testing.T) {
	a := assert.New(t)
	opt := &Options{}

	a.NotError(opt.setValue("pageSize", "25", true))
	a.Equal(opt.PageSize, 25)

	a.NotError(opt.setValue("pageSize", "45", false))
	a.Equal(opt.PageSize, 45)

	a.NotError(opt.setValue("commentsSize", "20", true))
	a.Equal(opt.CommentsSize, 20)

	a.Error(opt.setValue("commentsSize", "25", false))
}

func TestOptions_Get(t *testing.T) {
	a := assert.New(t)
	opt := &Options{PageSize: 22}

	val, found := opt.Get("pageSize")
	a.True(found).Equal(val, 22)
}

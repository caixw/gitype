// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/caixw/typing/models"
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

func TestOptions_UpdateFromOption(t *testing.T) {
	a := assert.New(t)
	opt := &Options{}

	o := &models.Option{Key: "pageSize", Group: "system", Value: "25"}
	a.NotError(opt.UpdateFromOption(o))
	a.Equal(opt.PageSize, 25)

	o = &models.Option{Key: "pageSize", Group: "system", Value: "45"}
	a.NotError(opt.UpdateFromOption(o))
	a.Equal(opt.PageSize, 45)
}

func TestOptions_GetValueByKey(t *testing.T) {
	a := assert.New(t)
	opt := &Options{PageSize: 22}

	val, found := opt.GetValueByKey("pageSize")
	a.True(found).Equal(val, 22)
}

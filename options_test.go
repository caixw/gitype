// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_fromMaps(t *testing.T) {
	a := assert.New(t)

	maps := []map[string]string{
		{
			"key":   "pageSize",
			"group": "system",
			"value": "50",
		},
		{ // 不存在的数据，会忽略
			"key":   "noexists",
			"group": "system",
			"value": "50",
		},
	}

	opt := &options{}
	a.NotError(opt.fromMaps(maps))
	a.Equal(opt.PageSize, 50)
}

func TestOptions_toMaps(t *testing.T) {
	a := assert.New(t)

	opt := &options{
		PageSize: 30,
		Pretty:   true,
	}

	maps, err := opt.toMaps()
	a.NotError(err)
	for _, item := range maps {
		if item["group"] == "system" && item["key"] == "pageSize" {
			a.Equal(item["value"], "30")
		}

		if item["group"] == "system" && item["key"] == "pretty" {
			a.Equal(item["value"], "true")
		}
	}
}

func TestOptions_updateFromOption(t *testing.T) {
	a := assert.New(t)
	opt := &options{}

	o := &option{Key: "pageSize", Group: "system", Value: "25"}
	a.NotError(opt.updateFromOption(o))
	a.Equal(opt.PageSize, 25)

	o = &option{Key: "pageSize", Group: "system", Value: "45"}
	a.NotError(opt.updateFromOption(o))
	a.Equal(opt.PageSize, 45)
}

func TestOptions_getValueByKey(t *testing.T) {
	a := assert.New(t)
	opt := &options{PageSize: 22}

	val, found := opt.getValueByKey("pageSize")
	a.True(found).Equal(val, 22)
}

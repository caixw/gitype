// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestConfig_replaceTitle(t *testing.T) {
	a := assert.New(t)
	conf := &config{
		Title: "TITLE",
	}

	test := func(title, value string) {
		a.Equal(value, conf.replaceTitle(title))
	}

	test("T%title%", "TTITLE")
	test("T", "T")
	test("", "")
	test("%title%", "TITLE")
	test("%Title%", "%Title%")
}

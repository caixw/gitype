// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSplitTags(t *testing.T) {
	a := assert.New(t)
	tags := []*Tag{
		{Slug: "1", Series: true},
		{Slug: "2", Series: false},
		{Slug: "3", Series: false},
		{Slug: "4", Series: true},
	}

	ts, series := splitTags(tags)
	a.Equal(len(ts), 2).Equal(len(series), 2)
}

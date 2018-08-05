// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTags(t *testing.T) {
	a := assert.New(t)

	tags, err := LoadTags(testdataPath)
	a.NotError(err).NotNil(tags)

	a.Equal(tags[0].Slug, "default1")
	a.Equal(tags[0].Color, "efefef")
	a.Equal(tags[0].Title, "默认1")
	a.Equal(tags[1].Slug, "default2")
}

func TestCheckTagsDup(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		{Slug: "1"},
		{Slug: "2"},
		{Slug: "3"},
	}
	a.NotError(checkTagsDup(tags))

	tags = append(tags, &Tag{Slug: "1"})
	a.Error(checkTagsDup(tags))
}

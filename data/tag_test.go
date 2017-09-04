// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestSplitTags(t *testing.T) {
	a := assert.New(t)
	tags := []*Tag{
		&Tag{Slug: "1", Series: true},
		&Tag{Slug: "2", Series: false},
		&Tag{Slug: "3", Series: false},
		&Tag{Slug: "4", Series: true},
	}

	ts, series := splitTags(tags)
	a.Equal(len(ts), 2).Equal(len(series), 2)
}

func TestLoadTags(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")

	tags, err := loadTags(p)
	a.NotError(err).NotNil(tags)

	a.Equal(tags[0].Slug, "default1")
	a.Equal(tags[0].Color, "efefef")
	a.Equal(tags[0].Title, "默认1")
	a.Equal(tags[1].Slug, "default2")
	a.Equal(tags[0].Permalink, vars.TagURL("default1", 0))
}

func TestCheckTagsDup(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		&Tag{Slug: "1"},
		&Tag{Slug: "2"},
		&Tag{Slug: "3"},
	}
	a.NotError(checkTagsDup(tags))

	tags = append(tags, &Tag{Slug: "1"})
	a.Error(checkTagsDup(tags))
}

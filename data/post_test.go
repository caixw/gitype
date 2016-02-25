// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestData_FindPost(t *testing.T) {
	a := assert.New(t)

	data := &Data{
		path: "./testdata/",
		Posts: []*Post{
			&Post{Slug: "default1"},
			&Post{Slug: "default2"},
		},
	}
	a.NotNil(data.FindPost("default1"))
	a.NotNil(data.FindPost("default2"))
	a.Nil(data.FindPost("default3"))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		&Tag{Slug: "default1"},
		&Tag{Slug: "default2"},
	}
	conf := &Config{
		Theme: "t1",
	}

	post, err := loadPost("./testdata/posts", "./testdata/posts/post1/meta.yaml", conf, tags)
	a.NotError(err).NotNil(post)
	a.Equal(len(post.Tags), 2).Equal(post.Tags[0].Slug, "default1")
	a.Equal(post.Modified, 0)
	a.Equal(post.Content, "<article>a1</article>\n")
}

func TestData_loadPosts(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		&Tag{Slug: "default1"},
		&Tag{Slug: "default2"},
	}
	conf := &Config{
		Theme: "t1",
	}

	d := &Data{
		path:   "./testdata",
		Tags:   tags,
		Config: conf,
	}
	a.NotError(d.loadPosts())
	a.Equal(len(d.Posts), 2)
	p2 := d.Posts[1]
	a.Equal(p2.Tags[0].Slug, "default1")
}

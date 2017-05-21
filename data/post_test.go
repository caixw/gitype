// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"path/filepath"
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		&Tag{Slug: "default1", Posts: make([]*Post, 0, 10)},
		&Tag{Slug: "default2", Posts: make([]*Post, 0, 10)},
	}
	conf := &Config{
		Theme: "t1",
	}

	post, err := loadPost(filepath.Clean("./testdata/data/posts"), filepath.Clean("./testdata/data/posts/post1/meta.yaml"), conf, tags)
	a.NotError(err).NotNil(post)
	a.Equal(len(post.Tags), 2).Equal(post.Tags[0].Slug, "default1")
	a.Equal(len(tags[0].Posts), 1) // 会同时增加标签的计数器
	a.Equal(post.Modified, 0)
	a.Equal(post.Template, "post") // 默认模板
	a.Equal(post.Content, "<article>a1</article>\n")

	post, err = loadPost(filepath.Clean("./testdata/data/posts"), filepath.Clean("./testdata/data/posts/folder/post2/meta.yaml"), conf, tags)
	a.NotError(err).NotNil(post)
	a.Equal(post.Slug, "folder/post2")
	a.Equal(post.Template, "t1") // 模板
}

func TestData_loadPosts(t *testing.T) {
	a := assert.New(t)

	tags := []*Tag{
		&Tag{Slug: "default1"},
		&Tag{Slug: "default2"},
	}
	urls := &URLS{
		Root:   "/root",
		Suffix: ".html",
		Post:   "posts",
	}
	conf := &Config{
		Theme: "t1",
		URLS:  urls,
	}

	d := &Data{
		path:   vars.NewPath("./testdata"),
		Tags:   tags,
		Config: conf,
	}
	a.NotError(d.loadPosts())
	a.Equal(len(d.Posts), 2)
	p2 := d.Posts[0]
	a.Equal(p2.Tags[0].Slug, "default1")
	a.Equal(p2.Permalink, "/root/posts/post1.html")
}

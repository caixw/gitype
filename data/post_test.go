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

func TestPost_sanitize(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")

	post, err := loadPost(p, filepath.Clean("./testdata/data/posts/post1/meta.yaml"))
	a.NotError(err).NotNil(post)

	a.NotError(post.sanitize())
	a.Equal(len(post.Tags), 0)     // 未调用 sanitize2 初始化
	a.Equal(post.Template, "post") // 默认模板
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")

	post, err := loadPost(p, filepath.Clean("./testdata/data/posts/post1/meta.yaml"))
	a.NotError(err).NotNil(post)
	a.Equal(len(post.Tags), 0) // 未调用 Data.sanitize2 初始化
	a.True(post.Modified.IsZero())
	a.Equal(post.Template, "") // 未调用 Post.sanitize
	a.Equal(post.Content, "<article>a1</article>\n")

	post, err = loadPost(p, filepath.Clean("./testdata/data/posts/folder/post2/meta.yaml"))
	a.NotError(err).NotNil(post)
	a.Equal(post.Slug, "folder/post2")
	a.Equal(post.Template, "t1") // 模板
}

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")

	posts, err := loadPosts(p)
	a.NotError(err).NotNil(posts)
	a.Equal(len(posts), 2)
}

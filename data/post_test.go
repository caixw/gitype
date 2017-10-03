// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/gitype/vars"
	"github.com/issue9/assert"
)

func TestCheckPostsDup(t *testing.T) {
	a := assert.New(t)

	posts := []*Post{
		{Slug: "1"},
		{Slug: "2"},
		{Slug: "3"},
	}
	a.NotError(checkPostsDup(posts))

	posts = append(posts, &Post{Slug: "1"})
	a.Error(checkPostsDup(posts))
}

func TestLoadPost(t *testing.T) {
	a := assert.New(t)

	post, err := loadPost(testdataPath, "/post1")
	a.NotError(err).NotNil(post)
	a.Equal(len(post.Tags), 0) // 未调用 Data.sanitize 初始化
	a.False(post.Modified.IsZero())
	a.Equal(post.Template, vars.PostTemplateName)
	a.Equal(post.Content, "<article>a1</article>\n")

	post, err = loadPost(testdataPath, "/folder/post2")
	a.NotError(err).NotNil(post)
	a.Equal(post.Slug, "/folder/post2")
	a.Equal(post.Template, "t1post") // 模板

	post, err = loadPost(testdataPath, "/draft")
	a.NotError(err).NotNil(post)
	a.True(post.Draft)
}

func TestLoadPosts(t *testing.T) {
	a := assert.New(t)

	posts, err := loadPosts(testdataPath)
	a.NotError(err).NotNil(posts)
	a.Equal(len(posts), 2) // 只有两条记录，Draft=true 的没有被加载
}

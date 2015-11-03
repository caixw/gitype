// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPostURL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/posts/1.html", PostURL(opt, "1"))
}

func TestTagURL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/tags/tag1.html", TagURL(opt, "tag1", 1))
	a.Equal("siteurl/tags/tag1.html?page=2", TagURL(opt, "tag1", 2))
}

// 生成文章列表url，首页不显示页码。
func TestPostsURL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/", PostsURL(opt, 1))
	a.Equal("siteurl/posts.html?page=2", PostsURL(opt, 2))
}

// 生成标签列表url，所有标签在一个页面显示，不分页。
func TestTagsURL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/tags.html", TagsURL(opt))
}

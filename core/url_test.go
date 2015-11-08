// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPostURL(t *testing.T) {
	Opt = &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/posts/1.html", PostURL("1"))
}

func TestTagURL(t *testing.T) {
	Opt = &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/tags/tag1.html", TagURL("tag1", 1))
	a.Equal("siteurl/tags/tag1.html?page=2", TagURL("tag1", 2))
}

// 生成文章列表url，首页不显示页码。
func TestPostsURL(t *testing.T) {
	Opt = &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/", PostsURL(1))
	a.Equal("siteurl/posts.html?page=2", PostsURL(2))
}

// 生成标签列表url，所有标签在一个页面显示，不分页。
func TestTagsURL(t *testing.T) {
	Opt = &Options{SiteURL: "siteurl/", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl/tags.html", TagsURL())
}

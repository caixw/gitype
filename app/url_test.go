// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_URL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("siteurl", opt.URL(""))
	a.Equal("siteurl/path/file.ext", opt.URL("/path/file.ext"))
	a.Equal("siteurl/posts/1.html", opt.URL(opt.PostURL("1")))
}

func TestOptions_PostURL(t *testing.T) {
	opt := &Options{SiteURL: "siteurl", Suffix: ".html"}
	a := assert.New(t)

	a.Equal("/posts/1.html", opt.PostURL("1")).
		Equal("siteurl/posts/1.html", opt.URL(opt.PostURL("1")))
}

func TestOptions_TagURL(t *testing.T) {
	opt := &Options{Suffix: ".html"}
	a := assert.New(t)

	a.Equal("/tags/tag1.html", opt.TagURL("tag1", 1))
	a.Equal("/tags/tag1.html?page=2", opt.TagURL("tag1", 2))
}

// 生成文章列表url，首页不显示页码。
func TestOptions_PostsURL(t *testing.T) {
	opt := &Options{Suffix: ".html"}
	a := assert.New(t)

	a.Equal("/posts.html", opt.PostsURL(1))
	a.Equal("/posts.html?page=2", opt.PostsURL(2))
}

// 生成标签列表url，所有标签在一个页面显示，不分页。
func TestOptions_TagsURL(t *testing.T) {
	opt := &Options{Suffix: ".html"}
	a := assert.New(t)

	a.Equal("/tags.html", opt.TagsURL())
}

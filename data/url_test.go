// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestData_URL(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{URL: "https://caixw.io"}}
	a.Equal(d.URL("/index.html"), "https://caixw.io/index.html")
}

func TestData_PostURL(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{Suffix: ".html"}}
	a.Equal(d.PostURL("1"), "/posts/1.html")
}

func TestData_PostsURL(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{Suffix: ".html"}}
	a.Equal(d.PostsURL(-1), "/")
	a.Equal(d.PostsURL(1), "/")
	a.Equal(d.PostsURL(2), "/posts.html?page=2")
}

func TestData_TagURL(t *testing.T) {
	a := assert.New(t)

	d := &Data{Config: &Config{Suffix: ".html"}}
	a.Equal(d.TagURL("tag1", -1), "/tags/tag1.html")
	a.Equal(d.TagURL("tag1", 1), "/tags/tag1.html")
	a.Equal(d.TagURL("tag1", 2), "/tags/tag1.html?page=2")
}

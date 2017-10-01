// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package url

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPost(t *testing.T) {
	a := assert.New(t)

	a.Equal(Post("1"), "/posts/1.html")
}

func TestPosts(t *testing.T) {
	a := assert.New(t)

	a.Equal(Posts(0), "/")
	a.Equal(Posts(1), "/")
	a.Equal(Posts(2), "/index.html?page=2")
}

func TestTag(t *testing.T) {
	a := assert.New(t)
	a.Equal(Tag("1", 0), "/tags/1.html")
	a.Equal(Tag("1", 1), "/tags/1.html")
	a.Equal(Tag("1", 2), "/tags/1.html?page=2")
}

func TestSearch(t *testing.T) {
	a := assert.New(t)

	a.Equal(Search("", 0), "/search.html")
	a.Equal(Search("", 1), "/search.html")
	a.Equal(Search("", 2), "/search.html?page=2")

	a.Equal(Search("q", 0), "/search.html?q=q")
	a.Equal(Search("q", 1), "/search.html?q=q")
	a.Equal(Search("q", 2), "/search.html?q=q&amp;page=2")
}

func TestThemes(t *testing.T) {
	a := assert.New(t)

	a.Equal(Theme(""), "/themes/")
	a.Equal(Theme("/"), "/themes/")
	a.Equal(Theme("/path"), "/themes/path")
	a.Equal(Theme("/path/1"), "/themes/path/1")
}

func TestAsset(t *testing.T) {
	a := assert.New(t)

	a.Equal(Asset("/"), "/posts/")
	a.Equal(Asset(""), "/posts/")
	a.Equal(Asset("/abc.png"), "/posts/abc.png")
	a.Equal(Asset("abc.png"), "/posts/abc.png")
}

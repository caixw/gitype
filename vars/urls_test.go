// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPostURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(PostURL("1"), "/posts/1.html")
}

func TestPostsURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(PostsURL(0), "/")
	a.Equal(PostsURL(1), "/")
	a.Equal(PostsURL(2), "/index.html?page=2")
}

func TestTagURL(t *testing.T) {
	a := assert.New(t)
	a.Equal(TagURL("1", 0), "/tags/1.html")
	a.Equal(TagURL("1", 1), "/tags/1.html")
	a.Equal(TagURL("1", 2), "/tags/1.html?page=2")
}

func TestSearchURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(SearchURL("", 0), "/search.html")
	a.Equal(SearchURL("", 1), "/search.html")
	a.Equal(SearchURL("", 2), "/search.html?page=2")

	a.Equal(SearchURL("q", 0), "/search.html?q=q")
	a.Equal(SearchURL("q", 1), "/search.html?q=q")
	a.Equal(SearchURL("q", 2), "/search.html?q=q&amp;page=2")
}

func TestThemesURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(ThemesURL(""), "/themes/")
	a.Equal(ThemesURL("/"), "/themes/")
	a.Equal(ThemesURL("/path"), "/themes/path")
	a.Equal(ThemesURL("/path/1"), "/themes/path/1")
}

func TestAssetURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(AssetURL("/"), "/posts/")
	a.Equal(AssetURL(""), "/posts/")
	a.Equal(AssetURL("/abc.png"), "/posts/abc.png")
	a.Equal(AssetURL("abc.png"), "/posts/abc.png")
}

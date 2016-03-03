// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"testing"

	"github.com/caixw/typing/data"
	"github.com/issue9/assert"
)

var (
	urlsTest = &data.URLS{
		Suffix: ".htm",
		Posts:  "/posts",
		Post:   "/post",
		Tags:   "/tags",
		Tag:    "/tag",
		Root:   "",
	}

	appTest = &app{data: &data.Data{URLS: urlsTest}}
)

func TestApp_postsURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(appTest.postsURL(1), "")
	a.Equal(appTest.postsURL(2), "/posts.htm?page=2")
}

func TestApp_postURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(appTest.postURL("1"), "/post/1.htm")
	a.Equal(appTest.postURL("/f/1"), "/post/f/1.htm")
}

func TestApp_tagsURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(appTest.tagsURL(), "/tags.htm")
}

func TestApp_tagURL(t *testing.T) {
	a := assert.New(t)

	a.Equal(appTest.tagURL("tag", 1), "/tag/tag.htm")
	a.Equal(appTest.tagURL("tag", 2), "/tag/tag.htm?page=2")
}

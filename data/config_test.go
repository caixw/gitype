// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestCheckRSS(t *testing.T) {
	a := assert.New(t)

	rss := &RSS{}
	a.Error(checkRSS("RSS", rss))

	// Size 错误
	rss.Size = 0
	a.Error(checkRSS("RSS", rss))
	rss.Size = -1
	a.Error(checkRSS("RSS", rss))

	// URL 错误
	rss.Size = 10
	a.Error(checkRSS("RSS", rss))

	rss.URL = "url"
	a.NotError(checkRSS("RSS", rss))
}

func TestCheckSitemap(t *testing.T) {
	a := assert.New(t)

	s := &Sitemap{}
	a.Error(checkSitemap(s))

	s.URL = "url"
	a.Error(checkSitemap(s))

	s.TagPriority = -1.0
	a.Error(checkSitemap(s))
	s.TagPriority = 1.1
	a.Error(checkSitemap(s))

	s.TagPriority = .8
	s.PostPriority = 0.9
	s.TagChangefreq = "never"
	s.PostChangefreq = "never"
	a.NotError(checkSitemap(s))
}

func TestIsChangereq(t *testing.T) {
	a := assert.New(t)

	a.False(isChangereq("n"))
	a.False(isChangereq(""))
	a.False(isChangereq("m"))

	a.True(isChangereq("never"))
}

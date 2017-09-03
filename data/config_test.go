// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestRSS_sanitize(t *testing.T) {
	a := assert.New(t)

	rss := &rssConfig{}
	conf := &config{
		Title: "title",
		RSS:   rss,
	}
	a.Error(rss.sanitize(conf, "rss"))

	// Size 错误
	rss.Size = 0
	a.Error(rss.sanitize(conf, "rss"))
	rss.Size = -1
	a.Error(rss.sanitize(conf, "rss"))

	// URL 错误
	rss.Size = 10
	a.Error(rss.sanitize(conf, "RSS"))

	rss.URL = "url"
	a.NotError(rss.sanitize(conf, "rss"))
	a.Equal(rss.Title, conf.Title)
}

func TestSitemap_sanitize(t *testing.T) {
	a := assert.New(t)

	s := &sitemapConfig{}
	a.Error(s.sanitize())

	s.URL = "url"
	a.Error(s.sanitize())

	s.Priority = -1.0
	a.Error(s.sanitize())
	s.Priority = 1.1
	a.Error(s.sanitize())

	s.Priority = .8
	s.PostPriority = 0.9
	s.Changefreq = "never"
	s.PostChangefreq = "never"
	a.NotError(s.sanitize())
	a.Equal(s.Type, vars.ContentTypeXML) // 默认值
}

func TestIsChangereq(t *testing.T) {
	a := assert.New(t)

	a.False(isChangereq("n"))
	a.False(isChangereq(""))
	a.False(isChangereq("m"))

	a.True(isChangereq("never"))
}

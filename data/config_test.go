// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	d := &Data{} // loadConfig 用到path.Data变量
	a.NotError(d.loadConfig("./testdata/meta/config.yaml")).NotNil(d.Config)
	conf := d.Config
	a.Equal(conf.Title, "title")
	a.Equal(conf.URL, "https://caixw.io")
	a.Equal(conf.Menus[0].URL, "url1")
	a.Equal(conf.Menus[1].Title, "title2")
}

func TestCheckConfig(t *testing.T) {
	a := assert.New(t)

	// PageSize
	conf := &Config{
		PageSize: -1,
	}
	a.Error(initConfig(conf))

	// LongDateFormat
	conf.LongDateFormat = "2006"
	a.Error(initConfig(conf))

	// ShortDateFormat
	conf.ShortDateFormat = "2006"
	a.Error(initConfig(conf))

	// UptimeFormat
	conf.UptimeFormat = "2006-01-02T17:01:22+0800"
	a.Error(initConfig(conf))

	// Author
	conf.PageSize = 1
	a.Error(initConfig(conf))

	// Author.Name
	conf.Author = &Author{}
	a.Error(initConfig(conf))

	// Title
	conf.Author.Name = "abc"
	a.Error(initConfig(conf))

	// URL
	conf.Title = "title"
	a.Error(initConfig(conf))
	// URL 格式错误
	conf.URL = "URL"
	a.Error(initConfig(conf))

	// themes
	conf.URL = "https://caixw.io"
	a.Error(initConfig(conf))

	// RSS
	conf.Theme = "t1"
	conf.RSS = &RSS{Title: "1", URL: "/", Size: 5}
	// conf.Atom = nil  // 当conf.Atom为nil时，不检测
	a.NotError(initConfig(conf))
}

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

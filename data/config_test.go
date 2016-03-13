// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	d := &Data{path: &vars.Path{Data: "./testdata"}} // loadConfig 用到path.Data变量
	a.NotError(d.loadConfig("./testdata/meta/config.yaml")).NotNil(d.Config)
	conf := d.Config
	a.Equal(conf.Title, "title")
	a.Equal(conf.URL, "https://caixw.io")
	a.Equal(conf.Menus[0].URL, "url1")
	a.Equal(conf.Menus[1].Title, "title2")
}

func TestFixedConfig(t *testing.T) {
	a := assert.New(t)

	conf := &Config{
		UptimeFormat: "",
		URL:          "https://caixw.io/",
	}
	a.NotError(fixedConfig(conf))
	a.Equal(conf.Uptime, 0).Equal(conf.URL, "https://caixw.io")
}

func TestCheckConfig(t *testing.T) {
	a := assert.New(t)

	// PageSize
	conf := &Config{
		PageSize: -1,
	}
	a.Error(checkConfig(conf, "./testdata"))

	// LongDateFormat
	conf.LongDateFormat = "2006"
	a.Error(checkConfig(conf, "./testdata"))

	// ShortDateFormat
	conf.ShortDateFormat = "2006"
	a.Error(checkConfig(conf, "./testdata"))

	// Author
	conf.PageSize = 1
	a.Error(checkConfig(conf, "./testdata"))

	// Author.Name
	conf.Author = &Author{}
	a.Error(checkConfig(conf, "./testdata"))

	// Title
	conf.Author.Name = "abc"
	a.Error(checkConfig(conf, "./testdata"))

	// URL
	conf.Title = "title"
	a.Error(checkConfig(conf, "./testdata"))
	// URL 格式错误
	conf.URL = "URL"
	a.Error(checkConfig(conf, "./testdata"))

	// themes
	conf.URL = "https://caixw.io"
	a.Error(checkConfig(conf, "./testdata"))

	// RSS.Title
	conf.Theme = "t1"
	conf.RSS = &RSS{Title: "1", URL: "/", Size: 5}
	// conf.Atom = nil  // 当conf.Atom为nil时，不检测
	a.NotError(checkConfig(conf, "./testdata"))
}

func TestCheckRSS(t *testing.T) {
	a := assert.New(t)

	rss := &RSS{}
	a.Error(checkRSS("RSS", rss))

	// Size 错误
	rss.Title = "title"
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

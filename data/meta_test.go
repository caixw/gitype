// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestData_loadTags(t *testing.T) {
	a := assert.New(t)

	data := &Data{
		path: vars.NewPath("./testdata"),
	}
	a.NotError(data.loadTags())
	a.NotNil(data.Tags)
	a.Equal(data.Tags[0].Slug, "default1")
	a.Equal(data.Tags[0].Color, "efefef")
	a.Equal(data.Tags[0].Title, "默认1")
	a.Equal(data.Tags[1].Slug, "default2")
	a.Equal(data.Tags[0].Permalink, "/tags/default1.html")

	t.Log(data.Tags[0])
}

func TestData_loadLinks(t *testing.T) {
	a := assert.New(t)

	data := &Data{path: vars.NewPath("./testdata")}
	a.NotError(data.loadLinks())
	a.True(len(data.Links) > 0)
	a.Equal(data.Links[0].Text, "text0")
	a.Equal(data.Links[0].URL, "url0")
	a.Equal(data.Links[1].Text, "text1")
	a.Equal(data.Links[1].URL, "url1")

	t.Log(data.Links)
}

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	d := &Data{path: vars.NewPath("./testdata")}
	a.NotError(d.loadConfig()).NotNil(d.Config)
	conf := d.Config
	a.Equal(conf.Title, "title")
	a.Equal(conf.URL, "https://caixw.io")
	a.Equal(conf.Menus[0].URL, "url1")
	a.Equal(conf.Menus[1].Title, "title2")
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

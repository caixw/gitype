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
	conf.RSS = &RSS{}
	a.Error(checkConfig(conf, "./testdata"))

	// RSS.Size
	conf.RSS.Title = "RSS"
	a.Error(checkConfig(conf, "./testdata"))
	conf.RSS.Size = -1
	a.Error(checkConfig(conf, "./testdata"))

	// Atom.Title
	conf.RSS.Size = 9
	conf.Atom = &RSS{}
	a.Error(checkConfig(conf, "./testdata"))

	// Atom.Size
	conf.Atom.Title = "RSS"
	a.Error(checkConfig(conf, "./testdata"))
	conf.Atom.Size = -1
	a.Error(checkConfig(conf, "./testdata"))

	// OK
	conf.Atom.Size = 100
	a.NotError(checkConfig(conf, "./testdata"))
}

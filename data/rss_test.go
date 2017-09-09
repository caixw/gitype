// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestRSSConfig_sanitize(t *testing.T) {
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

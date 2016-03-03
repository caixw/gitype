// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feeds"
	"github.com/caixw/typing/vars"
	"github.com/issue9/web"
)

type App struct {
	path    *vars.Path
	data    *data.Data
	updated int64

	// feed
	rssBuffer     *bytes.Buffer
	atomBuffer    *bytes.Buffer
	sitemapBuffer *bytes.Buffer
}

// 重新加载数据
func (a *App) reload() (err error) {
	a.data, err = data.Load(a.path)
	a.updated = time.Now().Unix()

	if a.data.Config.RSS != nil {
		a.rssBuffer, err = feeds.BuildRSS(a.data)
		if err != nil {
			return
		}
	}

	if a.data.Config.Atom != nil {
		a.atomBuffer, err = feeds.BuildAtom(a.data)
		if err != nil {
			return
		}
	}

	if a.data.Config.Sitemap != nil {
		a.sitemapBuffer, err = feeds.BuildSitemap(a.data)
		if err != nil {
			return
		}
	}
	return
}

func Run(p *vars.Path) error {
	a := &App{
		path: p,
	}

	// 加载程序配置
	data, err := ioutil.ReadFile(a.path.ConfApp)
	if err != nil {
		return err
	}
	conf := &web.Config{}
	if err = json.Unmarshal(data, conf); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		return err
	}

	// 初始化路由
	if err = a.initRoute(); err != nil {
		return err
	}

	return web.Run(conf)
}

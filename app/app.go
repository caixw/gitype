// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 核心处理模块，包括路由函数和页面渲染等。
// 会调用github.com/issue9/logs包的内容，调用之前需要初始化该包。
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

type app struct {
	path    *vars.Path
	module  *web.Module
	updated int64

	// 可重复加载的数据
	data          *data.Data
	rssBuffer     *bytes.Buffer
	atomBuffer    *bytes.Buffer
	sitemapBuffer *bytes.Buffer
}

// 重新加载数据
func (a *app) reload() (err error) {
	if a.data, err = data.Load(a.path); err != nil {
		return
	}

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

	a.updated = time.Now().Unix()

	// 重新初始化路由项
	return a.initRoute()
}

func Run(p *vars.Path) error {
	m, err := web.NewModule("front")
	if err != nil {
		return err
	}

	a := &app{
		path:   p,
		module: m,
	}

	// 加载数据
	if err = a.reload(); err != nil {
		return err
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
	return web.Run(conf)
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 核心处理模块，包括路由函数和页面渲染等。
// 会调用github.com/issue9/logs包的内容，调用之前需要初始化该包。
package app

import (
	"bytes"
	"html/template"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feeds"
	"github.com/caixw/typing/vars"
	"github.com/issue9/web"
)

type app struct {
	path     *vars.Path
	front    *web.Module
	conf     *config
	updated  int64 // 更新时间，一般为重新加载数据的时间
	adminTpl *template.Template

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

	conf, err := loadConfig(p.ConfApp)
	if err != nil {
		return err
	}

	a := &app{
		path:    p,
		front:   m,
		updated: time.Now().Unix(),
		conf:    conf,
	}

	// 初始化控制台相关操作
	if err := a.initAdmin(); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		return err
	}

	return web.Run(a.conf.Core)
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 核心处理模块，包括路由函数和页面渲染等。
// 会调用github.com/issue9/logs包的内容，调用之前需要初始化该包。
package app

import (
	"html/template"
	"net/http"
	"strconv"
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
	updated  int64  // 更新时间，一般为重新加载数据的时间
	etag     string // 所有页面都采用相同的etag
	adminTpl *template.Template
	data     *data.Data
}

// 重新加载数据
func (a *app) reload() error {
	data, err := data.Load(a.path)
	if err != nil {
		return err
	}
	a.data = data
	a.updated = time.Now().Unix()
	a.etag = strconv.FormatInt(a.updated, 10)
	a.front.Clean() // 清除路由项

	if err := a.initFrontRoute(); err != nil {
		return err
	}

	return a.initFeeds()
}

// 重新初始化路由项
func (a *app) initFrontRoute() error {
	urls := a.data.URLS
	p := a.front.Prefix(urls.Root)

	p.GetFunc(urls.Post+"/{slug:.+}"+urls.Suffix, a.pre(a.getPost)).
		GetFunc(vars.MediaURL+"/", a.pre(a.getMedia)).
		GetFunc(urls.Posts+urls.Suffix, a.pre(a.getPosts)).
		GetFunc(urls.Tag+"/{slug}"+urls.Suffix, a.pre(a.getTag)).
		GetFunc(urls.Tags+urls.Suffix+"{:.*}", a.pre(a.getTags)).
		GetFunc(urls.Themes+"/", a.pre(a.getThemes)).
		GetFunc("/", a.pre(a.getRaws))
	return nil
}

func (a *app) initFeeds() error {
	conf := a.data.Config
	p := a.front.Prefix(a.data.URLS.Root)

	if conf.RSS != nil {
		rss, err := feeds.BuildRSS(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.RSS.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(rss.Bytes())
		}))
	}

	if conf.Atom != nil {
		atom, err := feeds.BuildAtom(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Atom.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(atom.Bytes())
		}))
	}

	if conf.Sitemap != nil {
		sitemap, err := feeds.BuildSitemap(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Sitemap.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(sitemap.Bytes())
		}))
	}

	return nil
}

// 是否处于调试模式
func (a *app) isDebug() bool {
	return len(a.conf.Core.Pprof) > 0
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

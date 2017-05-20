// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 负责客户端的渲染
package client

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feeds"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

// Client 表示一个客户端渲染的相关集合
type Client struct {
	mux     *mux.Mux
	updated int64      // 更新时间，一般为重新加载数据的时间
	etag    string     // 所有页面都采用相同的 etag
	data    *data.Data // 加载的数据，每次加载都会被重置
	tpl     *template.Template
}

// New 声明一个新的 Client 实例
func New(datadir string, mux *mux.Mux) (*Client, error) {
	data, err := data.Load(datadir)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	return &Client{
		mux:     mux,
		data:    data,
		updated: now,
		etag:    strconv.FormatInt(now, 10),
	}, nil

	// init router
}

// 重新初始化路由项
func (a *app) initFrontRoute() error {
	urls := a.data.Config.URLS
	p := a.mux.Prefix(urls.Root)

	p.GetFunc(urls.Post+"/{slug}"+urls.Suffix, a.pre(a.getPost)).
		GetFunc(vars.MediaURL+"/*", a.pre(a.getMedia)).
		GetFunc(urls.Posts+urls.Suffix, a.pre(a.getPosts)).
		GetFunc(urls.Tag+"/{slug}"+urls.Suffix, a.pre(a.getTag)).
		GetFunc(urls.Tags+urls.Suffix+"{:.*}", a.pre(a.getTags)).
		GetFunc(urls.Themes+"/", a.pre(a.getThemes)).
		GetFunc(urls.Search+urls.Suffix+"{:.*}", a.pre(a.getSearch)).
		GetFunc("/", a.pre(a.getRaws))
	return nil
}

func (a *app) initFeeds() error {
	conf := a.data.Config
	p := a.mux.Prefix(a.data.Config.URLS.Root)

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

func Run(root string) error {
	logs.Info("程序工作路径为:", root)

	conf, err := loadConfig(filepath.Join(root, "conf", "app.json"))
	if err != nil {
		return err
	}

	a := &app{
		root:    root,
		mux:     mux.New(false, false, nil, nil),
		updated: time.Now().Unix(),
		conf:    conf,
	}

	// 初始化控制台相关操作
	if err := a.initAdmin(); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		logs.Error("app.Run:", err)
	}

	return http.ListenAndServeTLS(a.conf.Port, a.conf.CertFile, a.conf.KeyFile, a.mux)
}

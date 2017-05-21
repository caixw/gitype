// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 负责客户端的渲染
package client

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feeds"
	"github.com/issue9/handlers"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

// Client 表示一个客户端渲染的相关集合
type Client struct {
	mux     *mux.Mux
	debug   bool       // 是否为调试状态
	updated int64      // 更新时间，一般为重新加载数据的时间
	etag    string     // 所有页面都采用相同的 etag
	data    *data.Data // 加载的数据，每次加载都会被重置
	tpl     *template.Template
}

// New 声明一个新的 Client 实例
func New(datadir string, mux *mux.Mux, debug bool) (*Client, error) {
	d, err := data.Load(datadir)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	c := &Client{
		mux:     mux,
		debug:   debug,
		data:    d,
		updated: now,
		etag:    strconv.FormatInt(now, 10),
	}

	if err = c.initTemplate(); err != nil {
		return nil, err
	}

	if err = c.initRoutes(); err != nil {
		return nil, err
	}

	if err = c.initFeeds(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) initFeeds() error {
	conf := c.data.Config
	p := c.mux.Prefix(conf.URLS.Root)

	if conf.RSS != nil {
		rss, err := feeds.BuildRSS(c.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.RSS.URL, c.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(rss.Bytes())
		}))
	}

	if conf.Atom != nil {
		atom, err := feeds.BuildAtom(c.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Atom.URL, c.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(atom.Bytes())
		}))
	}

	if conf.Sitemap != nil {
		sitemap, err := feeds.BuildSitemap(c.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Sitemap.URL, c.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(sitemap.Bytes())
		}))
	}

	return nil
}

// 每次访问前需要做的预处理工作。
func (c *Client) pre(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.debug { // 调试状态，则每次都重新加载数据
			if err := c.Reload(); err != nil {
				logs.Error(err)
			}
		}

		// 输出访问日志
		logs.Infof("%v：%v", r.UserAgent(), r.URL)

		// 直接根据整个博客的最后更新时间来确认etag
		if r.Header.Get("If-None-Match") == c.etag {
			logs.Infof("304:%v", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", c.etag)
		handlers.CompressFunc(f).ServeHTTP(w, r)
	}
}

// Reload 重新加载数据
func (c *Client) Reload() error {
	// TODO

	return nil
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 负责客户端的渲染
package client

import (
	"net/http"

	"github.com/caixw/typing/buffer"
	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
)

// Client 展示给用户的前端页面。
type Client struct {
	mux    *mux.Mux
	path   *vars.Path
	buf    *buffer.Buffer
	routes []string // 记录路由项，释放时，需要删除这些路由项
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	b, err := buffer.New(path)
	if err != nil {
		return nil, err
	}

	c := &Client{
		mux:    mux,
		path:   path,
		buf:    b,
		routes: make([]string, 0, 10),
	}

	if err = c.initRoutes(); err != nil {
		return nil, err
	}

	c.initFeeds()

	return c, nil
}

// Free 释放所有的数据
func (c *Client) Free() {
	c.removeFeeds()
	c.removeRoutes()
}

func (c *Client) initFeeds() {
	conf := c.buf.Data.Config

	if conf.RSS != nil {
		c.mux.GetFunc(conf.RSS.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(c.buf.RSS)
		}))
	}

	if conf.Atom != nil {
		c.mux.GetFunc(conf.Atom.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(c.buf.Atom)
		}))
	}

	if conf.Sitemap != nil {
		c.mux.GetFunc(conf.Sitemap.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(c.buf.Sitemap)
		}))
	}
}

func (c *Client) removeFeeds() {
	conf := c.buf.Data.Config

	if conf.RSS != nil {
		c.mux.Remove(conf.RSS.URL)
	}

	if conf.Atom != nil {
		c.mux.Remove(conf.Atom.URL)
	}

	if conf.Sitemap != nil {
		c.mux.Remove(conf.Sitemap.URL)
	}
}

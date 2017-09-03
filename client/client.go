// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 客户路由处理
package client

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
)

// Client 处理用户请求
type Client struct {
	path     *vars.Path
	data     *data.Data
	mux      *mux.Mux
	patterns []string // 记录所有的路由项，方便释放时删除
	etag     string
	info     *info
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	client := &Client{
		path: path,
		mux:  mux,
		etag: strconv.FormatInt(d.Created.Unix(), 10),
		data: d,
	}
	client.info = client.newInfo()

	errFilter := func(fn func() error) {
		if err == nil {
			err = fn()
		}
	}

	// 依赖 data.Data 数据的相关操作
	errFilter(client.initRSS)
	errFilter(client.initAtom)
	errFilter(client.initSitemap)
	errFilter(client.initOpensearch)
	errFilter(client.initRoutes)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Created 返回当前数据的创建时间
func (client *Client) Created() time.Time {
	return client.data.Created
}

// Free 释放 Client 内容
func (client *Client) Free() {
	for _, pattern := range client.patterns {
		client.mux.Remove(pattern)
	}
}

func (client *Client) initOpensearch() error {
	if client.data.Opensearch == nil {
		return nil
	}

	o := client.data.Opensearch
	client.patterns = append(client.patterns, o.URL)
	client.mux.GetFunc(o.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, o.Type)
		w.Write(o.Content)
	}))

	return nil
}

func (client *Client) initSitemap() error {
	if client.data.Sitemap == nil {
		return nil
	}

	s := client.data.Sitemap
	client.patterns = append(client.patterns, s.URL)
	client.mux.GetFunc(s.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, s.Type)
		w.Write(s.Content)
	}))

	return nil
}

func (client *Client) initRSS() error {
	if client.data.RSS == nil {
		return nil
	}

	rss := client.data.RSS
	client.patterns = append(client.patterns, rss.URL)
	client.mux.GetFunc(rss.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, rss.Type)
		w.Write(rss.Content)
	}))

	return nil
}

func (client *Client) initAtom() error {
	if client.data.Atom == nil { // 不需要生成 atom
		return nil
	}

	atom := client.data.Atom
	client.patterns = append(client.patterns, atom.URL)
	client.mux.GetFunc(atom.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, atom.Type)
		w.Write(atom.Content)
	}))
	return nil
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 客户端路由处理
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

	client.addFeed(client.data.RSS)
	client.addFeed(client.data.Atom)
	client.addFeed(client.data.Sitemap)
	client.addFeed(client.data.Opensearch)

	if err := client.initRoutes(); err != nil {
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
		client.mux.Remove(pattern, http.MethodGet)
	}
}

func (client *Client) addFeed(feed *data.Feed) {
	if feed == nil {
		return
	}

	client.patterns = append(client.patterns, feed.URL)
	client.mux.GetFunc(feed.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, feed.Type)
		w.Write(feed.Content)
	}))
}

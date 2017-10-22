// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 对客户端请求的处理。
package client

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/path"
	"github.com/issue9/mux"
)

// Client 包含了整个可动态加载的数据以及路由的相关操作。
// 当需要重新加载数据时，只要获取一个新的 Client 实例即可。
type Client struct {
	path *path.Path
	mux  *mux.Mux

	data     *data.Data
	patterns []string // 记录所有的路由项，方便释放时删除
	etag     string
	info     *info
}

// New 声明一个新的 Client 实例
func New(path *path.Path, mux *mux.Mux) (*Client, error) {
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

	client.patterns = client.patterns[:0]
	client.Free()
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

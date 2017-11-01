// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 对客户端请求的处理。
package client

import (
	"net/http"
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
		data: d,
		info: newInfo(d),
	}

	return client, nil
}

// Mount 挂载路由
func (client *Client) Mount() error {
	if err := client.initFeedRoutes(); err != nil {
		return err
	}

	return client.initRoutes()
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

	// 释放 data 数据
	client.data.Free()
}

func (client *Client) initFeedRoutes() (err error) {
	handle := func(feed *data.Feed) {
		if err != nil || feed == nil {
			return
		}

		client.patterns = append(client.patterns, feed.URL)
		err = client.mux.HandleFunc(feed.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
			setContentType(w, feed.Type)
			w.Write(feed.Content)
		}), http.MethodGet)
	}

	handle(client.data.RSS)
	handle(client.data.Atom)
	handle(client.data.Sitemap)
	handle(client.data.Opensearch)

	return err
}

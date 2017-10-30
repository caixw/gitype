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
	"github.com/caixw/gitype/vars"
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
	updated  time.Time // 最后更新时间
	etag     string

	postsTicker     *time.Ticker
	postsTickerDone chan bool
}

// New 声明一个新的 Client 实例
func New(path *path.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	client := &Client{
		path:    path,
		mux:     mux,
		data:    d,
		updated: d.Created,
		etag:    vars.Etag(d.Created),
	}

	client.info = client.newInfo()

	client.addFeed(client.data.RSS)
	client.addFeed(client.data.Atom)
	client.addFeed(client.data.Sitemap)
	client.addFeed(client.data.Opensearch)

	if err := client.initRoutes(); err != nil {
		return nil, err
	}

	// 一切数据加载都没问题之后，开始运行更新服务。
	// 只有注册路由成功了，定时器开始工作才有意义。
	if d.Outdated != 0 {
		client.postsTicker = time.NewTicker(vars.OutdatedFrequency)
		client.postsTickerDone = make(chan bool, 1)
		client.runUpdateOutdatedServer()
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

	if client.postsTicker != nil {
		client.postsTicker.Stop()
		client.postsTickerDone <- true
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

func (client *Client) runUpdateOutdatedServer() {
	go func() {
		for {
			select {
			case <-client.postsTicker.C:
				client.updateOutdated()
			case <-client.postsTickerDone:
				return
			}
		}
	}()
}

func (client *Client) updateOutdated() {
	d := client.data

	if d.Outdated == 0 {
		return
	}

	now := d.CalcPostsOutdated()
	client.updated = now
	client.etag = vars.Etag(now)
}

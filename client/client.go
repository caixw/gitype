// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 是对 data 数据的再次加工以及所有非固定路由的处理，
// 方便重新加载数据时，可以直接释放整修 client 再重新生成。
package client

import (
	"html/template"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
)

// Client 处理用户请求
type Client struct {
	path       *vars.Path
	info       *info
	mux        *mux.Mux
	etag       string
	template   *template.Template // 主题编译后的模板
	rss        []byte
	atom       []byte
	sitemap    []byte
	opensearch []byte
	patterns   []string // 记录所有的路由项，方便翻译时删除

	Created int64 // 当前数据的加载时间
	Data    *data.Data
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	client := &Client{
		path:    path,
		mux:     mux,
		etag:    strconv.FormatInt(now, 10),
		Created: now,
		Data:    d,
	}
	client.info = client.newInfo()

	errFilter := func(fn func() error) {
		if err == nil {
			err = fn()
		}
	}

	// 依赖 data.Data 数据的相关操作
	errFilter(client.compileTemplate)
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

// Free 释放 Client 内容
func (client *Client) Free() {
	for _, pattern := range client.patterns {
		client.mux.Remove(pattern)
	}
}

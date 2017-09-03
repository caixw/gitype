// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 是对 data 数据的再次加工以及所有非固定路由的处理，
// 方便重新加载数据时，可以直接释放整个 client 再重新生成。
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
	path     *vars.Path
	data     *data.Data
	mux      *mux.Mux
	patterns []string // 记录所有的路由项，方便释放时删除
	etag     string

	// 由 data 延伸出的数据
	info       *info
	template   *template.Template // 主题编译后的模板
	rss        []byte
	atom       []byte
	sitemap    []byte
	opensearch []byte
	archives   []*archive

	Created time.Time // 当前数据的加载时间
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	client := &Client{
		path:    path,
		mux:     mux,
		etag:    strconv.FormatInt(now.Unix(), 10),
		Created: now,
		data:    d,
	}
	client.info = client.newInfo()

	errFilter := func(fn func() error) {
		if err == nil {
			err = fn()
		}
	}

	// 依赖 data.Data 数据的相关操作
	errFilter(client.compileTemplate)
	errFilter(client.initArchives)
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

func (client *Client) url(path string) string {
	return client.data.Config.URL + path
}

// Free 释放 Client 内容
func (client *Client) Free() {
	for _, pattern := range client.patterns {
		client.mux.Remove(pattern)
	}
}

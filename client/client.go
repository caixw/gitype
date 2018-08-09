// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 对客户端请求的处理。
package client

import (
	"net/http"
	"time"

	"github.com/issue9/logs"
	"github.com/issue9/middleware/compress"
	"github.com/issue9/mux"
	"github.com/issue9/web/encoding/html"
	"golang.org/x/text/message"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/path"
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
func New(path *path.Path) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	client := &Client{
		path: path,
		data: d,
		info: newInfo(d),
	}

	return client, nil
}

// Mount 挂载路由以及数据
func (client *Client) Mount(mux *mux.Mux, html *html.HTML) error {
	client.mux = mux

	html.SetTemplate(client.data.Theme.Template)

	// 为当前的语言注册一条数据
	// 使当前语言能被正确解析
	message.SetString(client.data.LanguageTag, "xx", "xx")

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

// 每次访问前需要做的预处理工作。
func (client *Client) prepare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs.Tracef("%s: %s", r.UserAgent(), r.URL) // 输出访问日志

		// 直接根据整个博客的最后更新时间来确认 etag
		if r.Header.Get("If-None-Match") == client.data.Etag {
			logs.Tracef("304: %s", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", client.data.Etag)
		compress.New(f, logs.ERROR(), map[string]compress.BuildCompressWriter{
			"gzip":    compress.NewGzip,
			"deflate": compress.NewDeflate,
		}).ServeHTTP(w, r)
	}
}

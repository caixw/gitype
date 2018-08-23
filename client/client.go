// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 对客户端请求的处理。
package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/issue9/logs"
	"github.com/issue9/middleware/compress"
	"github.com/issue9/mux"
	"github.com/issue9/utils"
	"github.com/issue9/web/context"
	"github.com/issue9/web/encoding"
	"github.com/issue9/web/encoding/html"
	"github.com/issue9/web/errorhandler"
	"golang.org/x/text/message"

	"github.com/caixw/gitype/client/page"
	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

// Client 包含了整个可动态加载的数据以及路由的相关操作。
// 当需要重新加载数据时，只要获取一个新的 Client 实例即可。
type Client struct {
	path *path.Path
	mux  *mux.Mux

	data     *data.Data
	patterns []string // 记录所有的路由项，方便释放时删除
	site     *page.Site
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
		site: page.NewSite(d),
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

	// 将所有的错误处理都指向同一个函数
	errorhandler.SetErrorHandler(client.renderError, 0)

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

// Page 生成页面
func (client *Client) page(typ string) *page.Page {
	return client.site.Page(typ, client.data)
}

// 输出当前内容到指定模板
func (client *Client) render(ctx *context.Context, p *page.Page, name string) {
	p.Charset = ctx.OutputCharsetName
	ctx.Render(http.StatusOK, html.Tpl(name, p), nil)
}

// 输出一个特定状态码下的错误页面。
// 若该页面模板不存在，则输出状态码对应的文本内容。
// 只查找当前主题目录下的相关文件。
// 只对状态码大于等于 400 的起作用。
func (client *Client) renderError(w http.ResponseWriter, code int) {
	logs.Debug("输出非正常状态码：", code)
	var data []byte

	// 根据情况输出内容，若不存在模板，则直接输出最简单的状态码对应的文本。
	filename := strconv.Itoa(code) + vars.TemplateExtension
	path := client.path.ThemesPath(client.data.Theme.ID, filename)
	if !utils.FileExists(path) {
		data = []byte(fmt.Sprintf("模板文件 %s 不存在\n", path))
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		data = []byte(err.Error())
	}

	w.Header().Set("Content-Type", errorContentType)
	w.WriteHeader(code)
	w.Write(data)
}

var errorContentType = encoding.BuildContentType("text/html", "utf-8")

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 核心处理模块，包括路由函数和页面渲染等。
// 会调用github.com/issue9/logs包的内容，调用之前需要初始化该包。
package app

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feeds"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/web"
)

type app struct {
	path     *vars.Path
	front    *web.Module        // 前台页面的模块
	conf     *config            // 配置内容
	updated  int64              // 更新时间，一般为重新加载数据的时间
	etag     string             // 所有页面都采用相同的etag
	adminTpl *template.Template // 后台管理的模板页面。
	data     *data.Data         // 加载的数据，每次加载都会被重置
}

// 重新加载数据
func (a *app) reload() error {
	data, err := data.Load(a.path)
	if err != nil {
		return err
	}
	a.data = data
	a.updated = time.Now().Unix()
	a.etag = strconv.FormatInt(a.updated, 10)
	a.front.Clean() // 清除路由项

	if err := a.initFrontRoute(); err != nil {
		return err
	}

	return a.initFeeds()
}

// 重新初始化路由项
func (a *app) initFrontRoute() error {
	urls := a.data.URLS
	p := a.front.Prefix(urls.Root)

	p.GetFunc(urls.Post+"/{slug:.+}"+urls.Suffix, a.pre(a.getPost)).
		GetFunc(vars.MediaURL+"/", a.pre(a.getMedia)).
		GetFunc(urls.Posts+urls.Suffix, a.pre(a.getPosts)).
		GetFunc(urls.Tag+"/{slug}"+urls.Suffix, a.pre(a.getTag)).
		GetFunc(urls.Tags+urls.Suffix+"{:.*}", a.pre(a.getTags)).
		GetFunc(urls.Themes+"/", a.pre(a.getThemes)).
		GetFunc(urls.Search+urls.Suffix+"{:.*}", a.pre(a.getSearch)).
		GetFunc("/", a.pre(a.getRaws))
	return nil
}

func (a *app) initFeeds() error {
	conf := a.data.Config
	p := a.front.Prefix(a.data.URLS.Root)

	if conf.RSS != nil {
		rss, err := feeds.BuildRSS(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.RSS.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(rss.Bytes())
		}))
	}

	if conf.Atom != nil {
		atom, err := feeds.BuildAtom(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Atom.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(atom.Bytes())
		}))
	}

	if conf.Sitemap != nil {
		sitemap, err := feeds.BuildSitemap(a.data)
		if err != nil {
			return err
		}

		p.GetFunc(conf.Sitemap.URL, a.pre(func(w http.ResponseWriter, r *http.Request) {
			w.Write(sitemap.Bytes())
		}))
	}

	return nil
}

func Run(path *vars.Path) error {
	logs.Info("程序工作路径为:", path.Root)

	front, err := web.NewModule("front")
	if err != nil {
		return err
	}

	conf, err := loadConfig(path.ConfApp)
	if err != nil {
		return err
	}

	a := &app{
		path:    path,
		front:   front,
		updated: time.Now().Unix(),
		conf:    conf,
	}

	// 初始化控制台相关操作
	if err := a.initAdmin(); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		logs.Error("app.Run:", err)
	}

	return web.Run(a.conf.Core)
}

// 输出一个特定状态码下的错误页面。若该页面模板不存在，则输出状态码对应的文本内容。
//
// 只查找当前主题目录下的相关文件。
// 只对状态码大于等于400的起作用。
func (a *app) renderStatusCode(w http.ResponseWriter, code int) {
	logs.Debug("输出非正常状态码：", code)
	if code < 400 {
		return
	}

	w.Header().Set("Content-Type", "text/html") // 需在 WriteHeader 之前设置
	w.WriteHeader(code)

	// 根据情况输出内容，若不存在模板，则直接输出最简单的状态码对应的文本。
	filename := strconv.Itoa(code) + ".html"
	path := filepath.Join(a.path.DataThemes, a.data.Config.Theme, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Errorf("page.renderStatusCode:查找模板文件[%v]时出现以下错误[%v]\n", path, err)

		// 若读取模板出错，则只输出与状态码相对应的基本文字描述
		data = []byte(http.StatusText(code))
	}
	w.Write(data)
}

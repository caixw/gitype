// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块，包括路由函数和页面渲染等。
// 会调用 github.com/issue9/logs 包的内容，调用之前需要初始化该包。
package app

import (
	"html/template"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"strings"

	"github.com/caixw/typing/buffer"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

const debugPprof = "/debug/pprof/"

type app struct {
	path     *vars.Path
	mux      *mux.Mux
	conf     *config
	buf      *buffer.Buffer
	adminTpl *template.Template // 后台管理的模板页面。
}

// Run 运行程序
func Run(path *vars.Path) error {
	logs.Info("程序工作路径为:", path.Root)

	conf, err := loadConfig(filepath.Join(path.ConfDir, "app.json"))
	if err != nil {
		return err
	}

	a := &app{
		path: path,
		mux:  mux.New(false, false, nil, nil),
		conf: conf,
	}

	// 初始化 webhooks
	a.mux.PostFunc(a.conf.WebhooksURL, a.postWebhooks)

	// 初始化控制台相关操作
	if err := a.initAdmin(); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		logs.Error(err)
	}

	// 路由由代码定义，不会更改，所以不需要在 a.reload() 中重新加载。
	if err = a.initRoutes(); err != nil {
		return err
	}

	h := a.buildHeader(a.buildPprof(a.mux))

	if !a.conf.HTTPS {
		return http.ListenAndServe(a.conf.Port, h)
	}

	go func() { // 对 80 端口的处理方式
		serveHTTP(a)
	}()
	return http.ListenAndServeTLS(a.conf.Port, a.conf.CertFile, a.conf.KeyFile, h)
}

func (a *app) buildHeader(h http.Handler) http.Handler {
	if len(a.conf.Headers) == 0 {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range a.conf.Headers {
			w.Header().Set(k, v)
		}
		h.ServeHTTP(w, r)
	})
}

// 根据 Config.Pprof 决定是否包装调试地址，调用前请确认是否已经开启 Pprof 选项
func (a *app) buildPprof(h http.Handler) http.Handler {
	if !a.conf.Pprof {
		return h
	}

	logs.Debug("开启了调试功能，地址为：", debugPprof)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, debugPprof) {
			h.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path[len(debugPprof):]
		switch path {
		case "cmdline":
			pprof.Cmdline(w, r)
		case "profile":
			pprof.Profile(w, r)
		case "symbol":
			pprof.Symbol(w, r)
		case "trace":
			pprof.Trace(w, r)
		default:
			pprof.Index(w, r)
		}
	}) // end return http.HandlerFunc
}

// 对 80 端口的处理方式
func serveHTTP(a *app) {
	switch a.conf.HTTPState {
	case "default":
		logs.Error(http.ListenAndServe(":80", a.mux))
	case "redirect":
		logs.Error(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 构建跳转链接
			url := r.URL
			url.Scheme = "HTTPS"
			url.Host = strings.Split(r.Host, ":")[0] + a.conf.Port

			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		})))
	case "disable":
		return
	}
}

// 重新加载数据
func (a *app) reload() error {
	// 移除 feed 路由
	if a.buf != nil {
		a.removeFeeds()
	}

	// 生成新的数据
	buf, err := buffer.New(a.path)
	if err != nil {
		return err
	}

	// 只有生成成功了，才替换老数据
	a.buf = buf

	// 重新生成 feed 路由
	a.initFeeds()

	return nil
}

func (a *app) initFeeds() {
	conf := a.buf.Data.Config

	if conf.RSS != nil {
		a.mux.GetFunc(conf.RSS.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.RSS)
		}))
	}

	if conf.Atom != nil {
		a.mux.GetFunc(conf.Atom.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.Atom)
		}))
	}

	if conf.Sitemap != nil {
		a.mux.GetFunc(conf.Sitemap.URL, a.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.buf.Sitemap)
		}))
	}
}

func (a *app) removeFeeds() {
	conf := a.buf.Data.Config

	if conf.RSS != nil {
		a.mux.Remove(conf.RSS.URL)
	}

	if conf.Atom != nil {
		a.mux.Remove(conf.Atom.URL)
	}

	if conf.Sitemap != nil {
		a.mux.Remove(conf.Sitemap.URL)
	}
}

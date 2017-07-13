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

type app struct {
	path     *vars.Path
	mux      *mux.Mux
	buf      *buffer.Buffer
	conf     *config            // 配置内容
	adminTpl *template.Template // 后台管理的模板页面。
}

// 重新加载数据
func (a *app) reload() error {
	buf, err := buffer.New(a.path)
	if err != nil {
		return err
	}
	a.buf = buf

	return nil
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
	if len(a.conf.Pprof) == 0 {
		return h
	}

	logs.Debug("开启了调试功能，地址为：", a.conf.Pprof)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, a.conf.Pprof) {
			h.ServeHTTP(w, r)
			return
		}

		path := r.URL.Path[len(a.conf.Pprof):]
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

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块。
//
// 包括 webhook 的处理以及整个 client 数据的替换等操作。
// 依赖 github.com/issue9/logs，确保该包已经被初始化。
package app

import (
	"html/template"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/caixw/typing/client"
	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

// 输出调试内容的地址，地址值固定，不能修改。
const debugPprof = "/debug/pprof/"

type app struct {
	path     *vars.Path
	mux      *mux.Mux
	conf     *config
	adminTpl *template.Template // 后台管理的模板页面。
	client   *client.Client
}

// Run 运行程序
func Run(path *vars.Path) error {
	logs.Info("程序工作路径为:", path.Root)

	conf := &config{}
	if err := helper.LoadYAMLFile(path.AppConfigFile, conf); err != nil {
		return err
	}
	if err := conf.sanitize(); err != nil {
		err.File = path.AppConfigFile
		return err
	}

	a := &app{
		path: path,
		mux:  mux.New(false, false, nil, nil),
		conf: conf,
	}

	// 初始化 webhooks
	a.mux.HandleFunc(a.conf.Webhook.URL, a.postWebhooks, a.conf.Webhook.Method)

	// 加载数据
	if err := a.reload(); err != nil {
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
		logs.Error(http.ListenAndServe(httpPort, a.mux))
	case "redirect":
		logs.Error(http.ListenAndServe(httpPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	// 释放旧数据
	if a.client != nil {
		a.client.Free()
	}

	// 生成新的数据
	c, err := client.New(a.path, a.mux)
	if err != nil {
		return err
	}

	// 只有生成成功了，才替换老数据
	a.client = c

	return nil
}

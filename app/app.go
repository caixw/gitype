// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块。
package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/caixw/gitype/client"
	"github.com/caixw/gitype/path"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

type app struct {
	path *path.Path
	mux  *mux.Mux
	conf *config

	// 当前是否正处在加载数据的状态，
	// 防止在 reload 一次调用未完成的情况下，再次调用 reload
	loading bool
	client  *client.Client
}

// Run 运行程序
func Run(path *path.Path, pprof, preview bool) error {
	logs.Info("程序工作路径为:", path.Root)

	conf, err := loadConfig(path)
	if err != nil {
		return err
	}

	a := &app{
		path: path,
		mux:  mux.New(false, false, nil, nil),
		conf: conf,
	}

	if preview {
		watcher, err := a.initWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		a.watch(watcher)
	} else {
		err = a.mux.HandleFunc(a.conf.Webhook.URL, a.postWebhooks, a.conf.Webhook.Method)
		if err != nil {
			return err
		}
	}

	// 加载数据，此时出错，只记录错误信息，但不中断执行
	if err = a.reload(); err != nil {
		logs.Error(err)
	}

	h := a.buildHandler(pprof)

	if !a.conf.HTTPS {
		return http.ListenAndServe(a.conf.Port, h)
	}

	go a.serveHTTP(h) // 对 80 端口的处理方式
	return http.ListenAndServeTLS(a.conf.Port, a.conf.CertFile, a.conf.KeyFile, h)
}

// 对 80 端口的处理方式
func (a *app) serveHTTP(h http.Handler) {
	switch a.conf.HTTPState {
	case httpStateDefault:
		logs.Error(http.ListenAndServe(httpPort, h))
	case httpStateRedirect:
		logs.Error(http.ListenAndServe(httpPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 构建跳转链接
			url := r.URL
			url.Scheme = "HTTPS"
			url.Host = strings.Split(r.Host, ":")[0] + a.conf.Port

			http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
		})))
	} // end switch
}

// 重新加载数据
func (a *app) reload() error {
	if a.loading {
		return errors.New("调用 reload 过于频繁")
	}

	a.loading = true
	defer func() { a.loading = false }()

	c, err := client.New(a.path, a.mux)
	if err != nil {
		return err
	}

	// 只有新数据生成成功了，才会释放旧数据，并加载新数据到路由中。
	if a.client != nil {
		a.client.Free()
	}
	a.client = c
	return a.client.Mount()
}

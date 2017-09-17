// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块。
//
// 包括 webhook 的处理以及整个 client 数据的替换等操作。
// 依赖 github.com/issue9/logs，确保该包已经被初始化。
package app

import (
	"net/http"
	"strings"

	"github.com/caixw/typing/client"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

// 输出调试内容的地址，地址值固定，不能修改。
const debugPprof = "/debug/pprof/"

type app struct {
	path   *vars.Path
	mux    *mux.Mux
	conf   *config
	client *client.Client
}

// Run 运行程序
func Run(path *vars.Path) error {
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

	// 初始化 webhooks
	err = a.mux.HandleFunc(a.conf.Webhook.URL, a.postWebhooks, a.conf.Webhook.Method)
	if err != nil {
		return err
	}

	// 加载数据，此时出错，只记录错误信息，但不中断执行
	if err = a.reload(); err != nil {
		logs.Error(err)
	}

	h := a.buildHandler()

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

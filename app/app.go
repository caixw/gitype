// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块，包括路由函数和页面渲染等。
// 会调用 github.com/issue9/logs 包的内容，调用之前需要初始化该包。
package app

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/caixw/typing/client"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

type app struct {
	path     *vars.Path
	mux      *mux.Mux
	conf     *config // 配置内容
	updated  int64   // 更新时间，一般为重新加载数据的时间
	client   *client.Client
	adminTpl *template.Template // 后台管理的模板页面。
}

// 重新加载数据
func (a *app) reload() error {
	if a.client != nil {
		a.client.Free()
	}

	c, err := client.New(a.path, a.mux)
	if err != nil {
		return err
	}
	a.client = c

	a.updated = time.Now().Unix()

	return nil
}

func Run(path *vars.Path) error {
	logs.Info("程序工作路径为:", path.Root)

	conf, err := loadConfig(filepath.Join(path.ConfDir, "app.json"))
	if err != nil {
		return err
	}

	a := &app{
		path:    path,
		mux:     mux.New(false, false, nil, nil),
		updated: time.Now().Unix(),
		conf:    conf,
	}

	// 初始化控制台相关操作
	if err := a.initAdmin(); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		logs.Error(err)
	}

	if a.conf.HTTPS {
		return http.ListenAndServeTLS(a.conf.Port, a.conf.CertFile, a.conf.KeyFile, a.mux)
	}
	return http.ListenAndServe(a.conf.Port, a.mux)
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package app 核心处理模块。
package app

import (
	"errors"
	"path/filepath"

	"github.com/issue9/logs"
	"github.com/issue9/mux"
	"github.com/issue9/web"
	"github.com/issue9/web/encoding"
	"github.com/issue9/web/encoding/html"

	"github.com/caixw/gitype/client"
	"github.com/caixw/gitype/path"
)

type app struct {
	path *path.Path

	// 当前是否正处在加载数据的状态，
	// 防止在 reload 一次调用未完成的情况下，再次调用 reload
	loading bool
	client  *client.Client
	html    *html.HTML
	webhook *webhook
	mux     *mux.Mux
}

// Run 运行程序
func Run(path *path.Path, preview bool) error {
	logs.Info("程序工作路径为:", path.Root)

	htmlMgr := html.New(nil)
	if err := encoding.AddMarshal("text/html", htmlMgr.Marshal); err != nil {
		return err
	}
	if err := encoding.AddMarshal("application/xhtml+xml", htmlMgr.Marshal); err != nil {
		return err
	}

	if err := web.Init(path.ConfDir); err != nil {
		return err
	}

	a := &app{
		path:    path,
		html:    htmlMgr,
		webhook: &webhook{},
		mux:     web.Mux(),
	}

	if err := web.LoadConfig(filepath.Join(path.ConfDir, "webhook.yaml"), a.webhook); err != nil {
		return err
	}

	if preview {
		watcher, err := a.initWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		a.watch(watcher)
	} else {
		conf := a.webhook
		if err := a.mux.HandleFunc(conf.URL, a.postWebhooks, conf.Method); err != nil {
			return err
		}
	}

	// 加载数据，此时出错，只记录错误信息，但不中断执行
	if err := a.reload(); err != nil {
		logs.Error(err)
	}

	return web.Serve()
}

// 重新加载数据
func (a *app) reload() error {
	if a.loading {
		return errors.New("调用 reload 过于频繁")
	}

	a.loading = true
	defer func() { a.loading = false }()

	c, err := client.New(a.path)
	if err != nil {
		return err
	}

	// 只有新数据生成成功了，才会释放旧数据，并加载新数据到路由中。
	if a.client != nil {
		a.client.Free()
	}
	a.client = c
	return a.client.Mount(a.mux, a.html)
}

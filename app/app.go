// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/path"
	"github.com/issue9/web"
)

const Version = "0.1.10.20160303" // 版本号

type App struct {
	path    *path.Path
	data    *data.Data
	updated int64
}

// 重新加载数据
func (a *App) reload() (err error) {
	a.data, err = data.Load(a.path)
	a.updated = time.Now().Unix()
	return
}

func Run(p *path.Path) error {
	a := &App{
		path: p,
	}

	// 加载程序配置
	data, err := ioutil.ReadFile(a.path.ConfApp)
	if err != nil {
		return err
	}
	conf := &web.Config{}
	if err = json.Unmarshal(data, conf); err != nil {
		return err
	}

	// 加载数据
	if err = a.reload(); err != nil {
		return err
	}

	// 初始化路由
	if err = a.initRoute(); err != nil {
		return err
	}

	web.Run(conf)
	return nil
}

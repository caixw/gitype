// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/core"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/install"
	"github.com/caixw/typing/themes"
	"github.com/issue9/mux"
	"github.com/issue9/web"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if install.Install() {
		return
	}

	err := core.Init()
	if err != nil {
		panic(err)
	}

	// themes
	if err = themes.Init(); err != nil {
		panic(err)
	}

	// admin
	admin.Init()

	// 初始化feed
	if err = feed.Init(); err != nil {
		panic(err)
	}

	if err := initModule(); err != nil {
		panic(err)
	}

	core.Cfg.Core.ErrHandler = mux.PrintDebug
	web.Run(core.Cfg.Core)
	core.DB.Close()
}

// 初始化模块，及与模块相对应的路由。
func initModule() error {
	// admin
	if err := admin.InitRoute(); err != nil {
		return err
	}

	// 初始化前台使用的api
	m, err := web.NewModule("front")
	if err != nil {
		return err
	}
	themes.InitRoute(m)

	feed.InitRoute(m)
	return nil
}

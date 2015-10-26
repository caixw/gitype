// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/core"
	"github.com/caixw/typing/install"
	"github.com/caixw/typing/sitemap"
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

	cfg, db, opt, err := core.Init()
	if err != nil {
		panic(err)
	}

	if err = themes.Init(cfg, opt, db); err != nil {
		panic(err)
	}

	// 初始化sitemap
	if err = sitemap.Init(cfg.TempDir + "sitemap.xml"); err != nil {
		panic(err)
	}

	// admin
	admin.Init(opt, db)

	if err := initModule(cfg); err != nil {
		panic(err)
	}

	cfg.Core.ErrHandler = mux.PrintDebug
	web.Run(cfg.Core)
	db.Close()
}

// 初始化模块，及与模块相对应的路由。
func initModule(cfg *core.Config) error {
	// admin
	m, err := web.NewModule("admin")
	if err != nil {
		return err
	}
	admin.InitRoute(m.Prefix(cfg.AdminAPIPrefix))

	// 初始化前台使用的api
	m, err = web.NewModule("front")
	if err != nil {
		return err
	}
	themes.InitRoute(m)

	//m.GetFunc("/rss", getRSS).
	//GetFunc("/rss/posts/{id}", getPostRSS)
	m.GetFunc("/sitemap.xml", sitemap.ServeHTTP)
	return nil
}

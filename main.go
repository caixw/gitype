// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/app"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/front"
	"github.com/issue9/logs"
	"github.com/issue9/web"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if install() {
		return
	}

	// app
	cfg, db, opt, stat, err := app.Init()
	if err != nil {
		panic(err)
	}

	// front
	if err = front.Init(cfg, db, opt, stat); err != nil {
		panic(err)
	}

	// admin
	if err := admin.Init(cfg, db, opt, stat); err != nil {
		panic(err)
	}

	// feed
	if err = feed.Init(db, opt); err != nil {
		panic(err)
	}

	web.Run(cfg.Core)
	db.Close()
	logs.Flush()
}

// 执行安装命令。
//
// 根据返回值来确定是否退出整个程序。
// 若返回true则表示当前已经执行完安装命令，可以退出整个程序，
// 否则表示当前程序没有从命令参数中获取安装指令，继续执行程序其它部分。
func install() bool {
	action := flag.String("init", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	flag.Parse()

	switch *action {
	case "config":
		if err := app.InstallConfig(); err != nil {
			panic(err)
		}

		return true
	case "db":
		if err := app.InstallDB(); err != nil {
			panic(err)
		}

		return true
	} // end switch

	return false
}

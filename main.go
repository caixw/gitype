// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/boot"
	"github.com/caixw/typing/core"
	"github.com/caixw/typing/feed"
	i "github.com/caixw/typing/install"
	"github.com/caixw/typing/themes"
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

	// boot
	cfg, db, err := boot.Init()
	if err != nil {
		panic(err)
	}

	// core
	opt, err := core.Init(db)
	if err != nil {
		panic(err)
	}

	// themes
	if err = themes.Init(cfg, db, opt); err != nil {
		panic(err)
	}

	// admin
	if err := admin.Init(cfg, db, opt); err != nil {
		panic(err)
	}

	// feed
	if err = feed.Init(cfg, db, opt); err != nil {
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
		if err := boot.Install(); err != nil {
			panic(err)
		}
		return true
	case "db":
		_, db, err := boot.Init()
		if err != nil {
			panic(err)
		}
		if err := i.Install(db); err != nil {
			panic(err)
		}
		return true
	} // end switch

	return false
}

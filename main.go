// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/app"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/front"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const usage = `typing 一个简单的博客程序，支持以下两个参数：
appdir 指定程序的数据存放路径，未指定，则为./config/；
install 若指定了值，则为相应的安装过程`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	appdir := flag.String("appdir", "./config/", "指定程序的数据存放目录")
	action := flag.String("install", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	flag.Parse()

	if len(*action) > 0 { // 执行安装过程
		if err := app.Install(*appdir, *action); err != nil {
			panic(err)
		}
	}

	// app
	if err := app.Init(*appdir); err != nil {
		panic(err)
	}

	// front
	if err := front.Init(); err != nil {
		panic(err)
	}

	// admin
	if err := admin.Init(); err != nil {
		panic(err)
	}

	// feed
	if err := feed.Init(); err != nil {
		panic(err)
	}

	app.Run()
	app.Close()
}

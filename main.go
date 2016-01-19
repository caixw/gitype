// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/app"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/front"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := app.Install("./config/"); err != nil {
		panic(err)
	}

	// app
	if err := app.Init("./config/"); err != nil {
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

// 执行安装命令。
//
// 根据返回值来确定是否退出整个程序。
// 若返回true则表示当前已经执行完安装命令，可以退出整个程序，
// 否则表示当前程序没有从命令参数中获取安装指令，继续执行程序其它部分。
func install(appdir string) bool {
	//action := flag.String("init", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	//flag.Parse()

	return true
}

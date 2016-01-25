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

const usage = `typing 一个简单博客程序。
源代码以MIT开源许可，并发布于github: https://github.com/caixw/typing

命令行语法：
 typing [options]

 options:
  -help    显示帮助信息；
  -appdir  指定程序的数据存放路径，未指定，则为./；
           若指定的路径不存在，在安装模式下会尝试创建；
  -install 若指定了值，则为执行相应的安装过程可选值为:
           -config 在appdir/config/下输出配置文件；
           -db     向数据库创建表及输出默认的数据项

常见用法：
 运行程序：    typing -appdir=/path
 输出配置文件: typing -appdir=/path -install=config
 创建数据库：  typing -appdir=/path -install=db`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	help := flag.Bool("help", false, "显示帮助信息")
	appdir := flag.String("appdir", "./", "指定程序的数据存放目录")
	action := flag.String("install", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if len(*action) > 0 { // 执行安装过程
		if err := app.Install(*appdir, *action); err != nil {
			panic(err)
		}
		return
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

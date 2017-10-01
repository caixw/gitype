// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 基于 Git 的博客系统。
package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/config"
	"github.com/caixw/typing/path"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
)

const usage = `%s 是一个基于 Git 的博客系统。
源代码以 MIT 开源许可发布于：%s


常见用法：

typing -pprof -appdir="./"
typing -appdir="./"


参数：

`

func main() {
	help := flag.Bool("h", false, "显示当前信息")
	version := flag.Bool("v", false, "显示程序的版本信息")
	pprof := flag.Bool("pprof", false, "是否在 /debug/pprof/ 启用调试功能")
	appdir := flag.String("appdir", "./", "指定运行的工作目录")
	init := flag.String("init", "", "初始化一个工作目录")
	flag.Usage = func() {
		fmt.Printf(usage, vars.Name, vars.URL)
		flag.PrintDefaults()
	}
	flag.Parse()

	switch {
	case *help:
		flag.Usage()
		return
	case *version:
		printVersion()
		return
	case len(*init) > 0:
		runInit(*init) // *init 指向的目录不存在时，会尝试创建
		return
	}

	path := path.New(*appdir)

	if err := logs.InitFromXMLFile(path.LogsConfigFile); err != nil {
		panic(err)
	}

	logs.Critical(app.Run(path, *pprof))
	logs.Flush()
}

func printVersion() {
	fmt.Printf("%s %s build with %s\n", vars.Name, vars.Version(), runtime.Version())

	if len(vars.CommitHash()) > 0 {
		fmt.Printf("Git commit hash:%s\n", vars.CommitHash())
	}
}

func runInit(root string) {
	if err := config.Init(path.New(root)); err != nil {
		panic(err)
	}

	fmt.Printf("操作成功，你现在可以在 %s 中修改具体的参数配置！\n", root)
}

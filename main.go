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
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
)

func main() {
	help := flag.Bool("h", false, "显示当前信息")
	version := flag.Bool("v", false, "显示程序的版本信息")
	appdir := flag.String("appdir", "./", "指定运行的工作目录")
	init := flag.String("init", "", "指定初始化的工作目录")
	flag.Usage = usage
	flag.Parse()

	switch {
	case *help:
		flag.Usage()
		return
	case *version:
		printVersion()
		return
	case len(*init) > 0:
		app.Init(vars.NewPath(*init))
		return
	}

	path := vars.NewPath(*appdir)

	err := logs.InitFromXMLFile(path.LogsConfigFile)
	if err != nil {
		panic(err)
	}

	logs.Critical(app.Run(path))
	logs.Flush()
}

func usage() {
	fmt.Fprintf(vars.CMDOutput, "%s 是一个基于 Git 的博客系统。\n", vars.AppName)
	fmt.Fprintf(vars.CMDOutput, "源代码以 MIT 开源许可发布于 Github: %s\n", vars.URL)

	fmt.Fprintln(vars.CMDOutput, "\n参数：")
	flag.CommandLine.SetOutput(vars.CMDOutput)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintf(vars.CMDOutput, "%s:%s build with %s\n", vars.AppName, vars.Version(), runtime.Version())
	if len(vars.CommitHash()) > 0 {
		fmt.Fprintf(vars.CMDOutput, "Git commit hash:%s\n", vars.CommitHash())
	}
}

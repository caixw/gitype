// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 简单的博客系统。
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
)

func main() {
	help := flag.Bool("h", false, "显示当前信息")
	version := flag.Bool("v", false, "显示程序的版本信息")
	appdir := flag.String("appdir", "./testdata", "指定运行的数据目录")
	flag.Usage = usage
	flag.Parse()

	switch {
	case *help:
		flag.Usage()
		return
	case *version:
		printVersion()
		return
	}

	path := vars.NewPath(*appdir)

	// 初始化日志
	err := logs.InitFromXMLFile(filepath.Join(path.ConfDir, "logs.xml"))
	if err != nil {
		panic(err)
	}

	logs.Critical(app.Run(path))
	logs.Flush()
}

func usage() {
	fmt.Fprintf(os.Stdout, "%v 一个简单博客程序。\n", vars.AppName)
	fmt.Fprintf(os.Stdout, "源代码以MIT开源许可，并发布于 Github: %v\n", vars.URL)

	fmt.Fprintln(os.Stdout, "\n参数:")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Fprintf(os.Stdout, "%v:%v build with %v\n", vars.AppName, vars.Version(), runtime.Version())
	if len(vars.CommitHash()) > 0 {
		fmt.Fprintf(os.Stdout, "git commit hash:%v\n", vars.CommitHash())
	}
}

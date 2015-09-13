// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	//"bufio"
	"io/ioutil"
	"os"
)

const (
	fileName    = "static.go" // 指定产生的文件名。
	packageName = "install"   // 指定包名。

	// 文件头部的警告内容
	warning = "// 该文件由make.go自动生成，请勿手动修改！\n\n"
)

var logFile = "./logs.xml"

func main() {
	w, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	w.WriteString(warning)

	// 输出包定义
	w.WriteString("package ")
	w.WriteString(packageName)
	w.WriteString("\n\n")

	w.WriteString("var LogFile=[]byte(`")
	data, err := ioutil.ReadFile(logFile)
	if err != nil {
		panic(err)
	}
	w.Write(data)
	w.WriteString("`)")
}

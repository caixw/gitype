// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"io/ioutil"
	"os"
)

const (
	file    = "./static.go" // 模板文件编译后保存的文件名
	pkgName = "static"      // 包的名称
)

// 模板文件名，及与其对应的可导出变量名
var templates = map[string]string{
	"./admin.html": "AdminHTML",
}

func compile(file *os.File, templateFile, varName string) error {
	src, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	file.WriteString("var ")
	file.WriteString(varName)
	file.WriteString(" = `")
	file.Write(src)
	file.WriteString("`")

	return nil
}

func main() {
	file, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("// 这是自动产生的文件，不需要修改")
	file.WriteString("\n\n")

	file.WriteString("package ")
	file.WriteString(pkgName)
	file.WriteString("\n\n")

	for filename, varName := range templates {
		if err := compile(file, filename, varName); err != nil {
			panic(err)
		}
	}
}

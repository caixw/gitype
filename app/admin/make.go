// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
)

const (
	file    = "./static.go" // 模板文件编译后保存的文件名
	pkgName = "admin"       // 包的名称
)

// 模板文件名，及与其对应的可导出变量名
var templates = map[string]string{
	"./admin.html": "AdminHTML",
}

func compile(buf *bytes.Buffer, templateFile, varName string) error {
	src, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	buf.WriteString("var ")
	buf.WriteString(varName)
	buf.WriteString(" = `")
	buf.Write(src)
	buf.WriteString("`")

	return nil
}

func main() {
	buf := bytes.NewBufferString("// 这是自动产生的文件，请不要修改！")
	buf.WriteString("\n\n")

	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteString("\n\n")

	for filename, varName := range templates {
		if err := compile(buf, filename, varName); err != nil {
			panic(err)
		}
	}

	// 格式化
	bs, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	file, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(bs)
	if err != nil {
		panic(err)
	}
}

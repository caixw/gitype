// Copyright 2018 by caixw, All rights reserved.
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
	packageName = "sw" // 需保持与当前包名相同
	fileheader  = "// 请勿修改此文件\n\n"
	output      = "./swjs.go"
	varName     = "swjs"
)

func main() {
	bs, err := ioutil.ReadFile("./sw.js")
	if err != nil {
		panic(err)
	}

	data := bytes.NewBufferString(fileheader)

	// 包名
	data.WriteString("package ")
	data.WriteString(packageName)
	data.Write([]byte{'\n'})

	// 变量内容
	data.WriteString("var ")
	data.WriteString(varName)
	data.WriteString("=[]byte(`")
	data.Write(bs)
	data.WriteString("`)")

	// 格式化
	bs, err = format.Source(data.Bytes())
	if err != nil {
		panic(err)
	}

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = file.Write(bs); err != nil {
		panic(err)
	}
}

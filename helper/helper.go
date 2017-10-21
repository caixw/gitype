// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package helper 一些通用的辅助类函数
package helper

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/caixw/gitype/vars"
	yaml "gopkg.in/yaml.v2"
)

// StatusError 标准的错误状态码输出函数，略作封装。
func StatusError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// LoadYAMLFile 加载 YAML 格式的文件 path 中的内容到 obj，obj 必须量个指针。
func LoadYAMLFile(path string, obj interface{}) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, obj)
}

// DumpYAMLFile 将 obj 转换成 YAML 格式的文本并输出到 path 指定的文件中
func DumpYAMLFile(path string, obj interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = file.Write(bs)
	return err
}

// DumpTextFile 将文本内容 text 输出到  path 指定的文件中
func DumpTextFile(path, text string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

// ReplaceContent 替换 content 中的 %content% 内容为 replacement
func ReplaceContent(content, replacement string) string {
	return strings.Replace(content, vars.ContentPlaceholder, replacement, -1)
}

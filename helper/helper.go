// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package helper 一些通用的辅助类函数
package helper

import (
	"io/ioutil"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// StatusError 标准的错误状态码输出函数，略作封装。
func StatusError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// LoadYAMLFile 加载 yaml 格式的文件 path 中的内容到 obj，obj 必须量个指针。
func LoadYAMLFile(path string, obj interface{}) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, obj)
}

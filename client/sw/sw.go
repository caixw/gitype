// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package sw 提供 service worker 的支持
package sw

import (
	"bytes"
)

// ServiceWorker SW 功能的管理
type ServiceWorker struct {
	caches map[string][]string
}

// New 声明新的 ServiceWorker 变量
func New() *ServiceWorker {
	return &ServiceWorker{
		caches: make(map[string][]string, 10),
	}
}

// Add 添加某一版本下的缓存文件路径。
func (sw *ServiceWorker) Add(ver string, paths ...string) {
	list, found := sw.caches[ver]
	if !found {
		sw.caches[ver] = paths
		return
	}

	list = append(list, paths...)
	sw.caches[ver] = list
}

// Bytes sw.js 的实际内容
func (sw *ServiceWorker) Bytes() []byte {
	content := new(bytes.Buffer)

	for ver, list := range sw.caches {
		content.WriteString("versions.set(")
		content.WriteString(ver)
		content.WriteString(", [")
		for _, item := range list {
			content.WriteByte('"')
			content.WriteString(item)
			content.WriteString("\",")
		}
		content.WriteString("]);")
	}

	return bytes.Replace(swjs, []byte("{{VERSIONS}}"), content.Bytes(), 1)
}

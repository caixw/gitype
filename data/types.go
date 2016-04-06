// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import "fmt"

// 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

// MetaError 用于描述加载data/meta下数据时的错误信息
type MetaError struct {
	File    string // 对应的文章名
	Field   string // 对应的字段名
	Message string // 对应的错误信息
}

func (e *MetaError) Error() string {
	return fmt.Sprintf("文件 %v 中的字段: %v 错误: %v", e.File, e.Field, e.Message)
}

// 排序接口
type posts []*Post

func (p posts) Less(i, j int) bool {
	switch {
	case p[i].Top && p[j].Top:
		return p[i].Created >= p[j].Created
	case p[i].Top:
		return false
	case p[j].Top:
		return true
	default:
		return p[i].Created >= p[j].Created
	}
}

func (p posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p posts) Len() int {
	return len(p)
}

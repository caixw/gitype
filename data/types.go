// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

// 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

// 描述链接内容
type Link struct {
	Icon  string `yaml:"icon,omitempty"`
	Title string `yaml:"title,omitempty"`
	URL   string `yaml:"url"`
	Text  string `yaml:"text'`
}

// 排序接口
type posts []*Post

func (p posts) Less(i, j int) bool {
	switch {
	case p[i].Top && p[j].Top:
		return p[i].Created < p[j].Created
	case p[i].Top:
		return true
	case p[j].Top:
		return false
	default:
		return p[i].Created < p[j].Created
	}
}

func (p posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p posts) Len() int {
	return len(p)
}

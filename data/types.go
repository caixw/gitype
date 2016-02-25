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

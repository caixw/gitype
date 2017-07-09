// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

// 与 URL 相关的一些定义，方便做一些自定义操作
const (
	Posts  = "/index"  // 列表页地址
	Post   = "/posts"  // 文章详细页地址
	Tags   = "/tags"   // 标签列表页地址
	Tag    = "/tags"   // 标签详细页地址
	Links  = "/links"  // 友情链接
	Search = "/search" // 搜索 URL，会加上 Suffix 作为后缀
	Themes = "/themes" // 主题地址
	Suffix = ".html"   // 地址后缀
)

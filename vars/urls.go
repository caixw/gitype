// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"path"
	"strconv"
)

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

func PostURL(slug string) string {
	return path.Join(Post, slug+Suffix)
}

func PostsURL(page int) string {
	if page <= 1 {
		return "/"
	}
	return Posts + Suffix + "?page=" + strconv.Itoa(page)
}

func TagURL(slug string, page int) string {
	url := path.Join(Tag, slug+Suffix)
	if page <= 1 {
		return url
	}

	return url + "?page=" + strconv.Itoa(page)
}

func SearchURL(q string, page int) string {
	url := Search + Suffix
	if len(q) > 0 {
		url += "?q=" + q
	}

	if page > 1 {
		if len(q) > 0 {
			url += "&"
		} else {
			url += "?"
		}
		url += "page=" + strconv.Itoa(page)
	}

	return url
}

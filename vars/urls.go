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
	posts  = "/index"  // 列表页地址
	post   = "/posts"  // 文章详细页地址
	tags   = "/tags"   // 标签列表页地址
	tag    = "/tags"   // 标签详细页地址
	links  = "/links"  // 友情链接
	search = "/search" // 搜索 URL，会加上 Suffix 作为后缀
	themes = "/themes" // 主题地址
	suffix = ".html"   // 地址后缀
)

// LinksURL 生成友情链接的 URL
func LinksURL() string {
	return links + suffix
}

// PostURL 构建文章的 URL
func PostURL(slug string) string {
	return path.Join(post, slug+suffix)
}

// PostsURL 构建文章列表的 URL
func PostsURL(page int) string {
	if page <= 1 {
		return "/"
	}
	return posts + suffix + "?page=" + strconv.Itoa(page)
}

// TagURL 构建标签的 URL
func TagURL(slug string, page int) string {
	url := path.Join(tag, slug+suffix)
	if page <= 1 {
		return url
	}

	return url + "?page=" + strconv.Itoa(page)
}

// TagsURL 生成标签列表的 URL
func TagsURL() string {
	return tags + suffix
}

// SearchURL 构建搜索页面的 URL
func SearchURL(q string, page int) string {
	url := search + suffix
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

// ThemesURL 构建主题文件 URL
func ThemesURL(p string) string {
	return path.Join(themes, p)
}

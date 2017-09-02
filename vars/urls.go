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
	root   = "/"     // 根地址
	suffix = ".html" // 地址后缀

	index    = root + "index"    // 列表页
	post     = root + "posts"    // 文章详细页
	tags     = root + "tags"     // 标签列表页
	tag      = root + "tags"     // 标签详细页
	links    = root + "links"    // 友情链接
	archives = root + "archives" // 归档
	search   = root + "search"   // 搜索
	themes   = root + "themes"   // 主题
)

var (
	linksURL    = links + suffix
	indexURL    = index + suffix
	tagsURL     = tags + suffix
	archivesURL = archives + suffix
	searchURL   = search + suffix
)

// LinksURL 生成友情链接的 URL
func LinksURL() string {
	return linksURL
}

// PostURL 构建文章的 URL
func PostURL(slug string) string {
	return path.Join(post, slug+suffix)
}

// PostsURL 构建文章列表的 URL
// 首页为返回 /
// 其它页面返回 /index.html?page=xx
func PostsURL(page int) string {
	if page <= 1 {
		return root
	}
	return indexURL + "?page=" + strconv.Itoa(page)
}

// IndexURL 构建索引首页的 URL
// 首页为返回 /index.html
// 其它页面返回 /index.html?page=xx
func IndexURL(page int) string {
	if page <= 1 {
		return indexURL
	}
	return indexURL + "?page=" + strconv.Itoa(page)
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
	return tagsURL
}

// ArchivesURL 生成归档页面的 URL
func ArchivesURL() string {
	return archivesURL
}

// SearchURL 构建搜索页面的 URL
func SearchURL(q string, page int) string {
	url := searchURL // 以下的 url+= 会改变 url 本身的值，所以不能直接使用 searchURL

	if len(q) > 0 {
		url += "?q=" + q
	}

	if page > 1 {
		if len(q) > 0 {
			url += "&amp;"
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

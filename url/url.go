// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package url

import (
	"path"
	"strconv"

	"github.com/caixw/typing/vars"
)

const (
	indexURL    = vars.URLRoot + vars.URLIndex + vars.URLSuffix    // 列表页
	postURL     = vars.URLRoot + vars.URLPost                      // 文章详细页
	tagsURL     = vars.URLRoot + vars.URLTags + vars.URLSuffix     // 标签列表页
	tagURL      = vars.URLRoot + vars.URLTag                       // 标签详细页
	linksURL    = vars.URLRoot + vars.URLLinks + vars.URLSuffix    // 友情链接
	archivesURL = vars.URLRoot + vars.URLArchives + vars.URLSuffix // 归档
	searchURL   = vars.URLRoot + vars.URLSearch + vars.URLSuffix   // 搜索
	themeURL    = vars.URLRoot + vars.URLTheme + "/"               // 主题目录前缀
	assetURL    = vars.URLRoot + vars.URLAsset + "/"               // 文章资源前缀
)

// Links 生成友情链接的 URL
func Links() string {
	return linksURL
}

// Post 构建文章的 URL
func Post(slug string) string {
	return path.Join(postURL, slug+vars.URLSuffix)
}

// Posts 构建文章列表的 URL
// 首页为返回 /
// 其它页面返回 /index.html?page=xx
func Posts(page int) string {
	if page <= 1 {
		return vars.URLRoot
	}
	return indexURL + "?" + vars.URLQueryPage + "=" + strconv.Itoa(page)
}

// Index 构建索引首页的 URL
// 首页为返回 /index.html
// 其它页面返回 /index.html?page=xx
func Index(page int) string {
	if page <= 1 {
		return indexURL
	}
	return indexURL + "?" + vars.URLQueryPage + "=" + strconv.Itoa(page)
}

// Tag 构建标签的 URL
func Tag(slug string, page int) string {
	url := path.Join(tagURL, slug+vars.URLSuffix)
	if page <= 1 {
		return url
	}

	return url + "?" + vars.URLQueryPage + "=" + strconv.Itoa(page)
}

// Tags 生成标签列表的 URL
func Tags() string {
	return tagsURL
}

// Archives 生成归档页面的 URL
func Archives() string {
	return archivesURL
}

// Search 构建搜索页面的 URL
func Search(q string, page int) string {
	url := searchURL // 以下的 url+= 会改变 url 本身的值，所以不能直接使用 searchURL

	if len(q) > 0 {
		url += "?" + vars.URLQueryQ + "=" + q
	}

	if page > 1 {
		if len(q) > 0 {
			url += "&amp;"
		} else {
			url += "?"
		}
		url += vars.URLQueryPage + "=" + strconv.Itoa(page)
	}

	return url
}

// Theme 构建主题文件 URL
func Theme(path string) string {
	return static(themeURL, path)
}

// Asset 构建一条用于指向资源的 URL
func Asset(path string) string {
	return static(assetURL, path)
}

func static(prefix, path string) string {
	if len(path) == 0 {
		return prefix
	}

	if path[0] == '/' {
		return prefix + path[1:]
	}
	return prefix + path
}

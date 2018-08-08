// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"path"
	"strconv"
)

// 查询参数名称的定义
const (
	URLQueryPage   = "page" // 查询参数 page
	URLQuerySearch = "q"    // 查询参数 q
)

// 与查询相关的一些自定义参数
//
// 上线之后请谨慎修改这些值，可能会让已经分离出去的链接变为无效链接。
//
// 用户可以通过查询参数按指定的格式进行精确查找，比如：
// title:abc 只查找标题中包含 abc 的文章，其中，title 关键字和分隔符 : 都可以自定义。
const (
	SearchKeySeparator = ':'
	SearchKeyTitle     = "title"
	SearchKeyTag       = "tag"
	SearchKeySeries    = "series"
)

// 与 URL 构成相关的配置项
const (
	urlSuffix = ".html" // 地址后缀

	indexURL    = "/index" + urlSuffix    // 列表页       /index.html
	postURL     = "/posts"                // 文章详细页   /posts
	tagsURL     = "/tags" + urlSuffix     // 标签列表页   /tags.html
	tagURL      = "/tags"                 // 标签详细页   /tags
	linksURL    = "/links" + urlSuffix    // 友情链接     /links.html
	archivesURL = "/archives" + urlSuffix // 归档         /archives.html
	searchURL   = "/search" + urlSuffix   // 搜索         /search.html
	themeURL    = "/themes/"              // 主题目录前缀 /themes/
	assetURL    = "/posts/"               // 文章资源前缀 /posts/
)

// LinksURL 生成友情链接的 URL
func LinksURL() string {
	return linksURL
}

// PostURL 构建文章的 URL
func PostURL(slug string) string {
	return path.Join(postURL, slug+urlSuffix)
}

// PostsURL 构建文章列表的 URL
// 首页为返回 /
// 其它页面返回 /index.html?page=xx
func PostsURL(page int) string {
	if page <= 1 {
		return "/"
	}
	return indexURL + "?" + URLQueryPage + "=" + strconv.Itoa(page)
}

// IndexURL 构建索引首页的 URL
// 首页为返回 /index.html
// 其它页面返回 /index.html?page=xx
func IndexURL(page int) string {
	if page <= 1 {
		return indexURL
	}
	return indexURL + "?" + URLQueryPage + "=" + strconv.Itoa(page)
}

// TagURL 构建标签的 URL
func TagURL(slug string, page int) string {
	url := path.Join(tagURL, slug+urlSuffix)
	if page <= 1 {
		return url
	}

	return url + "?" + URLQueryPage + "=" + strconv.Itoa(page)
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
		url += "?" + URLQuerySearch + "=" + q
	}

	if page > 1 {
		if len(q) > 0 {
			url += "&amp;"
		} else {
			url += "?"
		}
		url += URLQueryPage + "=" + strconv.Itoa(page)
	}

	return url
}

// ThemeURL 构建主题文件 URL
func ThemeURL(path string) string {
	return static(themeURL, path)
}

// AssetURL 构建一条用于指向资源的 URL
func AssetURL(path string) string {
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

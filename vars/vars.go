// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 定义一些全局变量、常量。相当于一个代码级别的配置内容。
package vars

import "time"

const (
	// Name 程序名称
	Name = "gitype"

	// URL 项目的地址
	URL = "https://github.com/caixw/gitype"

	// DateFormat 客户配置文件中所使用的的时间格式。
	// 所有的时间字符串，都将使用此格式去解析。
	//
	// 只负责时间的解析，如果是输出时间，则其格式由 meta/config.yaml 中定义。
	DateFormat = time.RFC3339

	// TemplateExtension 模板的扩展名
	TemplateExtension = ".html"

	// XMLIndentWidth XML 每一个 tab 的缩进量
	XMLIndentWidth = 4
)

// Cookie 的相关定义
const (
	CookieHTTPOnly = true
	CookieMaxAge   = 24 * 60 * 60

	// 主题名称在在传递过程中的名称
	CookieKeyTheme = "theme"
)

// 一些默认的字面文本内容。
const (
	NextPageText = "下一页"
	PrevPageText = "上一页"
)

// 与 URL 相关的一些定义
//
// 上线之后请谨慎修改这些值，会影响 URL 的路径结构。
const (
	URLRoot   = "/"     // 根地址
	URLSuffix = ".html" // 地址后缀

	URLIndex    = "index"    // 列表页
	URLPost     = "posts"    // 文章详细页
	URLTags     = "tags"     // 标签列表页
	URLTag      = "tags"     // 标签详细页
	URLLinks    = "links"    // 友情链接
	URLArchives = "archives" // 归档
	URLSearch   = "search"   // 搜索
	URLTheme    = "themes"   // 主题目录前缀
	URLAsset    = "posts"    // 文章资源前缀

	URLQueryPage = "page" // 查询参数 page
	URLQueryQ    = "q"    // 查询参数 q
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

// 目录名称的定义
const (
	DataDir = "data"
	ConfDir = "conf"

	PostsDir  = "posts"
	ThemesDir = "themes"
	MetaDir   = "meta"
	RawsDir   = "raws"
)

// 文件名的定义
const (
	AppConfigFilename  = "app.yaml"
	LogsConfigFilename = "logs.xml"

	ConfigFilename = "config.yaml"
	TagsFilename   = "tags.yaml"
	LinksFilename  = "links.yaml"

	PostMetaFilename    = "meta.yaml"
	PostContentFilename = "content.html"

	ThemeMetaFilename = "theme.yaml"
)

// 页面的类型，除了 PageIndex 其它的同时也是模板名称。
const (
	PageIndex    = "index" // 首页
	PagePosts    = "posts" // 除首页外的文章列表页
	PagePost     = "post"  // 同时也表示默认的模板名
	PageTags     = "tags"
	PageTag      = "tag"
	PageArchives = "archives"
	PageLinks    = "links"
	PageSearch   = "search"
)

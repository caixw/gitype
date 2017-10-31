// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 代码级别的配置内容。
//
// 所有可能需要修改的配置项以及算法都被集中到 vars 包中，
// 使用者可以根据自己需求随意修改此包以子包的内容。
package vars

import (
	"strconv"
	"time"
)

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

	// ContentPlaceholder 配置文件中的表示当前内容的占位符
	ContentPlaceholder = "%content%"

	// OutdatedFrequency outdated 的更新频率。
	// NOTE: 此值过小，有可能会影响服务器性能
	OutdatedFrequency = time.Hour * 24
)

// 目录名称的定义
const (
	DataFolderName = "data"
	ConfFolderName = "conf"

	PostsFolderName  = "posts"
	ThemesFolderName = "themes"
	MetaFolderName   = "meta"
	RawsFolderName   = "raws"
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

// Etag 根据一个时间，生成一段 Etag 字符串
func Etag(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

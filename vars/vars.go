// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 代码级别的配置内容。
//
// 所有可能需要修改的配置项以及算法都被集中到 vars 包中，
// 使用者可以根据自己需求随意修改此包以子包的内容。
//
// NOTE: vars 包的修改，可能会涉及到整个项目结构的改变，
// 比较适合用于一个新的博客内容，若是应用于已有的博客系统，
// 请确保这些修改不会对现有数据造成影响。
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
	PagePosts    = "posts" // 除首页外的文章列表页，与首页会有细微差别，比如标题
	PagePost     = "post"
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

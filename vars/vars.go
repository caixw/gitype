// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 定义一些全局变量、常量。相当于一个代码级别的配置内容。
package vars

import "time"

const (
	// Name 程序名称
	Name = "typing"

	// URL 项目的地址
	URL = "https://github.com/caixw/typing"

	// DateFormat 客户配置文件中所使用的的时间格式。
	// 所有的时间字符串，都将使用此格式去解析。
	//
	// 只负责时间的解析，如果是输出时间，则其格式由 meta/config.yaml 中定义。
	DateFormat = time.RFC3339

	// TemplateExtension 模板的扩展名
	TemplateExtension = ".html"

	// PostTemplateName 默认的文章模板名称
	PostTemplateName = "post"

	// ThemeName 主题名称在在传递过程中的名称
	ThemeName = "theme"
)

// 一些默认的字面文本内容。
const (
	NextPageText = "下一页"
	PrevPageText = "上一页"
)

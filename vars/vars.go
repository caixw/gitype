// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 定义一些全局变量、常量。
package vars

import "time"

const (
	// 主版本号，符合 semver 规范
	mainVersion = "0.10.6"

	// AppName 程序名称
	AppName = "typing"

	// URL 项目的地址
	URL = "https://github.com/caixw/typing"

	// DateFormat 客户配置文件中所使用的的时间格式。
	// 所有的时间字符串，都将使用此格式去解析。
	//
	// 只负责时间的解析，如果是输出时间，则其格式由 meta/config.yaml 中定义。
	DateFormat = time.RFC3339
)

// 默认的 mime type 类型。
//
// 一般情况下，用户无须修改此处内容，所有的 mime type 值都可以在
// data/meta/config.yaml 中手动指定，此处仅作为在其未指定的情况下
// 的一种默认值。
const (
	ContentTypeHTML       = "text/html"
	ContentTypeXML        = "application/xml"
	ContentTypeAtom       = "application/atom+xml"
	ContentTypeRSS        = "application/rss+xml"
	ContentTypeOpensearch = "application/opensearchdescription+xml"
)

var (
	buildDate  string
	commitHash string
	version    = mainVersion
)

func init() {
	if len(buildDate) > 0 {
		version += "+" + buildDate
	}
}

// Version 获取完整的版本号
func Version() string {
	return version
}

// CommitHash 获取最后一条代码提交记录的 hash 值。
func CommitHash() string {
	return commitHash
}

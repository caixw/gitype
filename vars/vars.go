// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package vars 定义一些全局变量、常量。
package vars

const (
	// 主版本号，符合 semver 规范
	mainVersion = "0.3.4"

	// AppName 程序名称
	AppName = "typing"

	// URL 项目的地址
	URL = "https://github.com/caixw/typing"

	// DateFormat 客户配置文件中所使用的的时间格式。
	// 所有的时间字符串，都将使用此格式去解析。
	//
	// 只负责时间的解析，如果是输出时间，则格式由 meta/config.yaml 中定义。
	DateFormat = "2006-01-02T15:04:05-0700"
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

// CommitHash 获取当前版本的 git commit hash
func CommitHash() string {
	return commitHash
}

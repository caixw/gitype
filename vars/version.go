// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

// 版本号
//
// major.minjor.patch，符合 semver 规范。
// 对外数据发生不兼容的时候，增加 major，其它值归零；
// 增加新功能是，增加 minjor，patch 归零；
// 修正 bug 则增加 patch；
// 其它普通改动，版本号不变。
//
// Version() 函数返回的值，并不能总是与 mainVersion 相同，
// 有可能还有编译日期等额外内容。
const mainVersion = "1.4.11"

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

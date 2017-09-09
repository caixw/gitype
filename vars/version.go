// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

// 主版本号，符合 semver 规范
const mainVersion = "1.4.8"

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

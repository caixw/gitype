// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import "path/filepath"

// 一些文件名的定义
const (
	appConfigFilename  = "app.json"
	logsConfigFilename = "logs.xml"
)

// 一些目录名称的定义
const (
	dataDir = "data"
	confDir = "conf"

	postsDir  = "posts"
	themesDir = "themes" // NOTE: 此值需要与 themes 保持一致
	metaDir   = "meta"
	rawsDir   = "raws"
)

// Path 表示的文件路径信息
type Path struct {
	Root string // 项目的根目录，即 -appdir 参数指定的目录

	ConfDir string // 项目下的配置文件所在目录
	DataDir string // 项目下数据文件所在的目录，即 Git 数据所在的目录

	// 数据目录下的子目录
	PostsDir  string
	ThemesDir string
	MetaDir   string
	RawsDir   string

	AppConfigFile  string
	LogsConfigFile string
}

// NewPath 声明一个新的 Path
func NewPath(root string) *Path {
	dataDir := filepath.Join(root, dataDir)
	confDir := filepath.Join(root, confDir)

	p := &Path{
		Root: root,

		DataDir: dataDir,
		ConfDir: confDir,

		PostsDir:  filepath.Join(dataDir, postsDir),
		ThemesDir: filepath.Join(dataDir, themesDir),
		MetaDir:   filepath.Join(dataDir, metaDir),
		RawsDir:   filepath.Join(dataDir, rawsDir),
	}

	p.AppConfigFile = p.ConfPath(appConfigFilename)
	p.LogsConfigFile = p.ConfPath(logsConfigFilename)

	return p
}

// MetaPath 获取 data/meta/ 下的文件
func (p *Path) MetaPath(file string) string {
	return filepath.Join(p.MetaDir, file)
}

// ConfPath 获取 conf/ 下的文件
func (p *Path) ConfPath(file string) string {
	return filepath.Join(p.ConfDir, file)
}

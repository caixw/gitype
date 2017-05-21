// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import "path/filepath"

const (
	dataDir = "data"
	confDir = "conf"

	postsDir  = "posts"
	themesDir = "themes"
	metaDir   = "meta"
	rawsDir   = "raws"
)

// Path 表示的文件路径信息
type Path struct {
	Root string

	DataDir string
	ConfDir string

	PostsDir  string
	ThemesDir string
	MetaDir   string
	RawsDir   string
}

// NewPath 声明一个新的 Path
func NewPath(root string) *Path {
	dataDir := filepath.Join(root, dataDir)
	confDir := filepath.Join(root, confDir)

	return &Path{
		Root: root,

		DataDir: dataDir,
		ConfDir: confDir,

		PostsDir:  filepath.Join(dataDir, postsDir),
		ThemesDir: filepath.Join(dataDir, themesDir),
		MetaDir:   filepath.Join(dataDir, metaDir),
		RawsDir:   filepath.Join(dataDir, rawsDir),
	}
}

// MetaPath 获取 data/meta/ 下的文件
func (p *Path) MetaPath(file string) string {
	return filepath.Join(p.MetaDir, file)
}

// ConfPath 获取 conf/ 下的文件
func (p *Path) ConfPath(file string) string {
	return filepath.Join(p.ConfDir, file)
}

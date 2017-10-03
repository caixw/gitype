// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package path 处理程序的路径信息
package path

import (
	"path/filepath"

	"github.com/caixw/gitype/vars"
)

// Path 表示 gitype 的目录结构信息。
//
// gitype 拥有一个固定的目录结构，程序根据这个目录结构加载相关的数据信息，
// Path 可以在指定根目录的情况下，预先生成所有的目录结构路径，方便其它地方调用。
type Path struct {
	Root string // 项目的根目录

	ConfDir string // 项目下的配置文件所在目录
	DataDir string // 项目下数据文件所在的目录，即 Git 数据所在的目录

	// 数据目录下的子目录
	PostsDir  string
	ThemesDir string
	MetaDir   string
	RawsDir   string

	AppConfigFile  string
	LogsConfigFile string

	MetaConfigFile string
	MetaLinksFile  string
	MetaTagsFile   string
}

// New 声明一个新的 Path
func New(root string) *Path {
	dataDir := filepath.Join(root, vars.DataDir)
	confDir := filepath.Join(root, vars.ConfDir)

	p := &Path{
		Root: root,

		DataDir: dataDir,
		ConfDir: confDir,

		PostsDir:  filepath.Join(dataDir, vars.PostsDir),
		ThemesDir: filepath.Join(dataDir, vars.ThemesDir),
		MetaDir:   filepath.Join(dataDir, vars.MetaDir),
		RawsDir:   filepath.Join(dataDir, vars.RawsDir),
	}

	p.AppConfigFile = p.ConfPath(vars.AppConfigFilename)
	p.LogsConfigFile = p.ConfPath(vars.LogsConfigFilename)

	p.MetaConfigFile = p.MetaPath(vars.ConfigFilename)
	p.MetaLinksFile = p.MetaPath(vars.LinksFilename)
	p.MetaTagsFile = p.MetaPath(vars.TagsFilename)

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

// RawsPath 获取 raws/ 下的文件
func (p *Path) RawsPath(file string) string {
	return filepath.Join(p.RawsDir, file)
}

// ThemeMetaPath 返回指定主题下的描述文件
func (p *Path) ThemeMetaPath(theme string) string {
	return filepath.Join(p.ThemesDir, theme, vars.ThemeMetaFilename)
}

// PostPath 返回某一篇文章下的文件名
func (p *Path) PostPath(slug, filename string) string {
	slug = filepath.FromSlash(slug) // slug 有可能带路径分隔符
	return filepath.Join(p.PostsDir, slug, filename)
}

// PostMetaPath 返回某一篇文章下的 meta.yaml 文件地址
func (p *Path) PostMetaPath(slug string) string {
	return p.PostPath(slug, vars.PostMetaFilename)
}

// PostContentPath 返回某一篇文章下的文章内容的文件地址
func (p *Path) PostContentPath(slug string) string {
	return p.PostPath(slug, vars.PostContentFilename)
}

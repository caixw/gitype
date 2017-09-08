// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import "path/filepath"

// 一些文件名的定义
const (
	appConfigFilename  = "app.yaml"
	logsConfigFilename = "logs.xml"

	configFilename = "config.yaml"
	tagsFilename   = "tags.yaml"
	linksFilename  = "links.yaml"

	PostMetaFilename    = "meta.yaml"
	postContentFilename = "content.yaml"

	themeMetaFilename = "theme.yaml"
)

// 一些目录名称的定义
const (
	dataDir = "data"
	confDir = "conf"

	postsDir  = "posts"
	themesDir = "themes"
	metaDir   = "meta"
	rawsDir   = "raws"
)

// Path 表示文件的路径信息
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

	MetaConfigFile string
	MetaLinksFile  string
	MetaTagsFile   string
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

	p.MetaConfigFile = p.MetaPath(configFilename)
	p.MetaLinksFile = p.MetaPath(linksFilename)
	p.MetaTagsFile = p.MetaPath(tagsFilename)

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

// RawsPath 获取 rasw/ 下的文件
func (p *Path) RawsPath(file string) string {
	return filepath.Join(p.RawsDir, file)
}

// ThemeMetaPath 返回指定主题下的描述文件
func (p *Path) ThemeMetaPath(theme string) string {
	return filepath.Join(p.ThemesDir, theme, themeMetaFilename)
}

// PostPath 返回某一篇文章下的文件名
func (p *Path) PostPath(slug, filename string) string {
	slug = filepath.FromSlash(slug) // slug 有可能带路径分隔符
	return filepath.Join(p.PostsDir, slug, filename)
}

// PostMetaPath 返回某一篇文章下的 Meta.yaml 文件地址
func (p *Path) PostMetaPath(slug string) string {
	return p.PostPath(slug, PostMetaFilename)
}

// PostContentPath 返回某一篇文章下的文章内容的文件地址
func (p *Path) PostContentPath(slug string, filename string) string {
	if len(filename) == 0 {
		filename = postContentFilename
	}
	return p.PostPath(slug, filename)
}

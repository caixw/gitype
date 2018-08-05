// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/caixw/gitype/helper"
	p "github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
	"github.com/issue9/utils"
)

var defaultRobots = `User-agent:*
Disallow:/themes/`

var defaultPostContent = `<section>about
</section>`

var defaultTheme = &Theme{
	Name:        "default",
	Version:     "0.1.0",
	Description: "默认主题",
	URL:         vars.URL,
}

var defaultPostMeta = &Post{
	Title:    "about",
	Tags:     "default",
	Created:  time.Now().Format(vars.DateFormat),
	Modified: time.Now().Format(vars.DateFormat),
	State:    StateLast,
}

var defaultConfig = &Config{
	Title:           "Title",
	Language:        language,
	Subtitle:        "subtitle",
	PageSize:        20,
	LongDateFormat:  "2006-01-02 15:04:05",
	ShortDateFormat: "2006-01-02",
	Type:            contentTypeHTML,
	Author: &Author{
		Name: vars.Name,
		URL:  vars.URL,
	},
	License: &Link{
		Title: "署名 4.0 国际 (CC BY 4.0)",
		URL:   "https://creativecommons.org/licenses/by/4.0/deed.zh",
	},

	Theme:        "default",
	UptimeFormat: time.Now().Format(vars.DateFormat),
	Archive: &ArchiveConfig{
		Type:   ArchiveTypeYear,
		Format: "2006 年",
	},

	RSS: &RSSConfig{
		Title: "RSS",
		URL:   "/rss.xml",
		Type:  contentTypeRSS,
		Size:  20,
	},

	Pages: map[string]*Page{
		vars.PageArchives: &Page{
			Title:    archivesTitle,
			Keywords: "存档,归档,archive,archives",
		},
	},
}

var defaultLinks = []*Link{
	{
		Text: vars.Name,
		URL:  vars.URL,
	},
	{
		Text: "caixw",
		URL:  "https://caixw.io",
	},
}

var defaultTags = []*Tag{
	{
		Title: "默认",
		Slug:  "default",
	},
}

// Init 初始化 data 下的基本数据结构
func Init(path *p.Path) error {
	fmt.Println(path.DataDir)
	if !utils.FileExists(path.DataDir) {
		if err := os.Mkdir(path.DataDir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := initRaws(path); err != nil {
		return err
	}

	if err := initMeta(path); err != nil {
		return err
	}

	if err := initPosts(path); err != nil {
		return err
	}

	return initThemes(path)
}

// 初始化 data/meta 目录下的数据
func initMeta(path *p.Path) error {
	if !utils.FileExists(path.MetaDir) {
		if err := os.Mkdir(path.MetaDir, os.ModePerm); err != nil {
			return err
		}
	}

	// data/meta/config.yaml
	if err := helper.DumpYAMLFile(path.MetaConfigFile, defaultConfig); err != nil {
		return err
	}

	// data/meta/links.yaml
	if err := helper.DumpYAMLFile(path.MetaLinksFile, defaultLinks); err != nil {
		return err
	}

	// data/meta/tags.yaml
	return helper.DumpYAMLFile(path.MetaTagsFile, defaultTags)
}

// 初始化 data/raws 目录下的数据
func initRaws(path *p.Path) error {
	if !utils.FileExists(path.RawsDir) {
		if err := os.Mkdir(path.RawsDir, os.ModePerm); err != nil {
			return err
		}
	}

	// robots.txt
	return helper.DumpTextFile(path.RawsPath("robots.txt"), defaultRobots)
}

// 初始化 data/posts 目录下数据
func initPosts(p *p.Path) error {
	slug := path.Join(strconv.Itoa(time.Now().Year()), "about")

	dir := filepath.Join(p.PostsDir, slug)
	if !utils.FileExists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := helper.DumpYAMLFile(p.PostMetaPath(slug), defaultPostMeta); err != nil {
		return err
	}

	// content.html
	return helper.DumpTextFile(p.PostContentPath(slug), defaultPostContent)
}

// 初始化 data/themes 目录
func initThemes(path *p.Path) error {
	if !utils.FileExists(path.ThemesDir) {
		if err := os.Mkdir(path.ThemesDir, os.ModePerm); err != nil {
			return err
		}
	}

	return helper.DumpYAMLFile(path.ThemeMetaPath("default"), defaultTheme)
}

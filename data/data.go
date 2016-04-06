// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 负责加data目录下的数据。
// 会调用github.com/issue9/logs包的内容，调用之前需要初始化该包。
package data

import (
	"html/template"

	"github.com/caixw/typing/vars"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	Config   *Config            // 配置内容
	URLS     *URLS              // 自定义URL
	Tags     []*Tag             // map对顺序是未定的，所以使用slice
	Links    []*Link            // 友情链接
	Template *template.Template // 当前主题模板
	Posts    []*Post            // 所有的文章列表
}

// Load 函数用于加载一份新的数据。
func Load(path *vars.Path) (*Data, error) {
	d := &Data{}

	if err := d.loadMeta(path); err != nil {
		return nil, err
	}

	// 加载文章
	if err := d.loadPosts(path.DataPosts); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Data) loadMeta(path *vars.Path) error {
	// urls
	if err := d.loadURLS(path.DataURLS); err != nil {
		return err
	}

	// tags
	if err := d.loadTags(path.DataTags); err != nil {
		return err
	}

	// links
	if err := d.loadLinks(path.DataLinks); err != nil {
		return err
	}

	// config
	if err := d.loadConfig(path.DataConf); err != nil {
		return err
	}

	// theme
	themes, err := getThemesName(path.DataThemes)
	if err != nil {
		return err
	}
	found := false
	for _, theme := range themes {
		if theme == d.Config.Theme {
			found = true
			break
		}
	}
	if !found {
		return &MetaError{File: "config.yaml", Message: "该主题并不存在", Field: "Theme"}
	}

	// 加载主题的模板
	return d.loadTemplate(path.DataThemes)
}

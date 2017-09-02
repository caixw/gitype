// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/caixw/typing/vars"
)

// Theme 表示主题信息
type Theme struct {
	ID          string  `yaml:"-"`           // 主题的唯一 ID
	Name        string  `yaml:"name"`        // 主题名称
	URL         string  `yaml:"url"`         // 网站
	Version     string  `yaml:"version"`     // 主题的版本号
	Description string  `yaml:"description"` // 主题的描述信息
	Author      *Author `yaml:"author"`      // 作者
	Path        string  `yaml:"-"`           // 主题所在的目录
}

func loadThemes(path *vars.Path) ([]*Theme, error) {
	dir := path.ThemesDir
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(fs) == 0 {
		return nil, errors.New("未找到任何主题文件")
	}

	themes := make([]*Theme, 0, len(fs))

	for _, file := range fs {
		if !file.IsDir() {
			continue
		}
		theme, err := loadTheme(path, file.Name())
		if err != nil {
			return nil, err
		}
		themes = append(themes, theme)
	}

	sort.SliceStable(themes, func(i, j int) bool {
		return themes[i].Name >= themes[j].Name
	})

	return themes, nil
}

// id 主题当前目录名称
func loadTheme(path *vars.Path, id string) (*Theme, error) {
	p := path.ThemeMetaPath(id)

	theme := &Theme{}
	if err := loadYamlFile(p, theme); err != nil {
		return nil, err
	}

	theme.Path = filepath.Dir(p)
	theme.ID = id

	return theme, nil
}

func (theme *Theme) sanitize() *FieldError {
	if len(theme.Name) == 0 {
		return &FieldError{Message: "不能为空", Field: "name"}
	}

	if theme.Author != nil {
		if err := theme.Author.sanitize(); err != nil {
			return err
		}
	}

	return nil
}

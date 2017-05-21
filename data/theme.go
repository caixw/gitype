// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

func (d *Data) loadThemes() error {
	dir := d.path.ThemesDir

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return errors.New("未找到任何主题文件")
	}

	d.Themes = make([]*Theme, 0, len(fs))

	for _, file := range fs {
		if !file.IsDir() {
			continue
		}
		theme, err := loadTheme(dir, file.Name())
		if err != nil {
			return err
		}
		d.Themes = append(d.Themes, theme)
	}

	sort.SliceStable(d.Themes, func(i, j int) bool {
		switch {
		case d.Themes[i].Actived:
			return true
		case d.Themes[j].Actived:
			return true
		default:
			return d.Themes[i].Name >= d.Themes[j].Name
		}
	})

	return nil
}

// dir 主题所在的目录
// id 主题当前目录名称
func loadTheme(dir, id string) (*Theme, error) {
	path := filepath.Join(dir, id, "theme.yaml")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	theme := &Theme{}
	if err = yaml.Unmarshal(data, theme); err != nil {
		return nil, fmt.Errorf("解板[%v]出错:%v", path, err)
	}

	if len(theme.Name) == 0 {
		return nil, &FieldError{File: path, Message: "不能为空", Field: "name"}
	}
	if theme.Author != nil {
		// err 必须是一个新变量，否则判断会一直是 true
		if err := theme.Author.check(); err != nil {
			return nil, err
		}
	}

	theme.Path = path
	theme.ID = id

	return theme, nil
}

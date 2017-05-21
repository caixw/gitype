// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	yaml "gopkg.in/yaml.v2"
)

func (d *Data) loadThemes() error {
	dir := filepath.Join(d.Root, ThemesDir)

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return errors.New("未找到任何主题文件")
	}

	themes := make([]*Theme, 0, len(fs))

	for _, file := range fs {
		if !file.IsDir() {
			continue
		}
		theme, err := loadTheme(dir, file.Name())
		if err != nil {
			return err
		}
		themes = append(themes, theme)
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

// 加载主题目录下的所有主题。
// path 主题所在的目录。
func getThemesName(path string) ([]string, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(fs) == 0 {
		return nil, errors.New("未找到任何主题文件")
	}

	themes := make([]string, 0, len(fs))

	for _, file := range fs {
		if !file.IsDir() {
			continue
		}

		themes = append(themes, file.Name())
	}

	return themes, nil
}

// 加载主题的模板。
// dir 模板所在的目录。
func (d *Data) loadTemplate(dir string) error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"ldate":    d.longDateFormat,
		"sdate":    d.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return path.Join(d.Config.URLS.Themes, p) },
	}

	var err error
	d.Template, err = template.New("").
		Funcs(funcMap).
		ParseGlob(filepath.Join(dir, d.Config.Theme, "*.html"))
	return err
}

func rfc3339DateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(time.RFC3339)
}

// 根据options中的格式显示长日期
func (d *Data) longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(d.Config.LongDateFormat)
}

// 根据options中的格式显示短日期
func (d *Data) shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(d.Config.ShortDateFormat)
}

// 将内容显示为html内容
func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签。
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}

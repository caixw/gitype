// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"html/template"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

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

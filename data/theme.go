// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/caixw/typing/vars"
)

// 肯定存在的模板名称，检测模板是否存在时，会用到。
// 该值会被 checkTemplates 改变，不能用于其它地方。
var templates = []string{
	vars.DefaultPostTemplateName,
	"posts",
	"tags",
	"tag",
	"links",
	"archives",
	"search",
}

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

	if len(theme.Name) == 0 {
		return nil, &FieldError{File: theme.Path, Message: "不能为空", Field: "name"}
	}

	if theme.Author != nil {
		if err := theme.Author.sanitize(); err != nil {
			err.Field = theme.Path
			return nil, err
		}
	}

	return theme, nil
}

// 编译主题的模板。
func (d *Data) compileTemplate() error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"unix":     unix,
		"ldate":    d.longDateFormat,
		"sdate":    d.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return vars.ThemesURL(p) },
	}

	tpl, err := template.New("client").
		Funcs(funcMap).
		ParseGlob(filepath.Join(d.Theme.Path, "*"+vars.TemplateExtension))
	if err != nil {
		return err
	}
	d.Template = tpl

	return d.checkTemplatesExists()
}

// 检测模板名称是否在模板中真实存在
func (d *Data) checkTemplatesExists() error {
	for _, post := range d.Posts {
		for _, tpl := range templates {
			if tpl != post.Template {
				templates = append(templates, post.Template)
			}
		}
	}

	// 模板定义未必是按文件分的，所以不能简单地判断文件是否存在
	for _, tpl := range templates {
		if nil != d.Template.Lookup(tpl) {
			continue
		}

		return fmt.Errorf("模板 %s 未定义", tpl)
	}

	return nil
}

func rfc3339DateFormat(t time.Time) interface{} {
	return t.Format(time.RFC3339)
}

// 转换成 unix 时间戳
func unix(t time.Time) interface{} {
	return t.Unix()
}

func (d *Data) longDateFormat(t time.Time) interface{} {
	return t.Format(d.Config.LongDateFormat)
}

func (d *Data) shortDateFormat(t time.Time) interface{} {
	return t.Format(d.Config.ShortDateFormat)
}

// 将内容显示为 HTML 内容
func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签。
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}

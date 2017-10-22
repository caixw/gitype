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

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
	"github.com/caixw/gitype/vars/url"
)

// Theme 表示主题信息
type Theme struct {
	ID          string             `yaml:"-"`    // 唯一 ID，即当前目录名称
	Name        string             `yaml:"name"` // 名称，不必唯一，可以与 ID 值不同。
	Path        string             `yaml:"-"`    // 主题目录，绝对路径
	Version     string             `yaml:"version"`
	Description string             `yaml:"description"`
	URL         string             `yaml:"url,omitempty"`
	Author      *Author            `yaml:"author"`
	Template    *template.Template `yaml:"-"` // 当前主题的预编译结果
}

func loadThemes(path *path.Path) ([]*Theme, error) {
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

	return themes, nil
}

// id 主题当前目录名称
func loadTheme(path *path.Path, id string) (*Theme, error) {
	p := path.ThemeMetaPath(id)

	theme := &Theme{}
	if err := helper.LoadYAMLFile(p, theme); err != nil {
		return nil, err
	}

	theme.Path = filepath.Dir(p)
	theme.ID = id

	if len(theme.Name) == 0 {
		return nil, &helper.FieldError{File: path.ThemeMetaPath(theme.ID), Message: "不能为空", Field: "name"}
	}

	if theme.Author != nil {
		if err := theme.Author.sanitize(); err != nil {
			err.Field = path.ThemeMetaPath(theme.ID)
			return nil, err
		}
	}

	return theme, nil
}

func (d *Data) sanitizeThemes(conf *config) error {
	var defaultTheme *Theme
	for _, theme := range d.Themes { // 检测配置文件中的主题是否存在
		if theme.ID == conf.Theme {
			defaultTheme = theme
			break
		}
	}

	if defaultTheme == nil {
		return &helper.FieldError{File: d.path.MetaConfigFile, Message: "该主题并不存在", Field: "theme"}
	}

	sort.SliceStable(d.Themes, func(i, j int) bool {
		// 确保默认主题在第一个位置
		if defaultTheme == d.Themes[i] {
			return true
		}
		if defaultTheme == d.Themes[j] {
			return false
		}

		return d.Themes[i].Name < d.Themes[j].Name
	})

	return d.compileTemplates()
}

// 编译主题的模板。
func (d *Data) compileTemplates() error {
	templates := d.templatesName()

	snippets, err := d.snippetsTemplate()
	if err != nil {
		return err
	}

	// 编译各个主题
	for _, theme := range d.Themes {
		theme.Template, err = snippets.Clone()
		if err != nil {
			return err
		}

		_, err = theme.Template.ParseGlob(filepath.Join(theme.Path, "*"+vars.TemplateExtension))
		if err != nil {
			return err
		}

		// 检测模板名称是否在模板中真实存在
		// 模板定义未必是按文件分的，所以不能简单地判断文件是否存在
		for _, tpl := range templates {
			if nil == theme.Template.Lookup(tpl) {
				return fmt.Errorf("模板 %s 未定义", tpl)
			}
		}
	}

	return nil
}

// 获取公用的代码片段模板
func (d *Data) snippetsTemplate() (*template.Template, error) {
	funs := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"unix":     unix,
		"ldate":    d.longDate,
		"sdate":    d.shortDate,
		"rfc3339":  rfc3339Date,
		"themeURL": func(p string) string { return url.Theme(p) },
	}

	return template.New("snippets").
		Funcs(funs).
		ParseGlob(filepath.Join(d.path.ThemesDir, "*"+vars.TemplateExtension))
}

// 获取所有的模板名称，除了固定的模板名称之外，
// 文章可以自定义模板名称。
func (d *Data) templatesName() []string {
	var templates = []string{
		vars.PagePost,
		vars.PagePosts,
		vars.PageTags,
		vars.PageTag,
		vars.PageLinks,
		vars.PageArchives,
		vars.PageSearch,
	}

	// 只有文章页可以自定义模板名称
	for _, post := range d.Posts {
		// 默认模板名，肯定已存在于 templates 变量中
		if post.Template == vars.PagePost {
			continue
		}

		for _, tpl := range templates {
			if tpl != post.Template {
				templates = append(templates, post.Template)
			}
		}
	}

	return templates
}

func rfc3339Date(t time.Time) interface{} {
	return t.Format(time.RFC3339)
}

// 转换成 unix 时间戳
func unix(t time.Time) interface{} {
	return t.Unix()
}

func (d *Data) longDate(t time.Time) interface{} {
	return t.Format(d.longDateFormat)
}

func (d *Data) shortDate(t time.Time) interface{} {
	return t.Format(d.shortDateFormat)
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

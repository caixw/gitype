// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

// Theme 表示主题信息
type Theme struct {
	ID          string  `yaml:"-"`    // 唯一 ID，即当前目录名称
	Name        string  `yaml:"name"` // 名称，不必唯一，可以与 ID 值不同。
	Version     string  `yaml:"version"`
	Description string  `yaml:"description"`
	URL         string  `yaml:"url,omitempty"`
	Author      *Author `yaml:"author"`

	template        *template.Template // 当前主题的预编译结果
	longDateFormat  string             // 长时间的显示格式
	shortDateFormat string             // 短时间的显示格式
}

// 查找与 conf.Theme 相同的主题。若找不到，则返回 errors
func findTheme(path *path.Path, conf *config) (*Theme, error) {
	dir := path.ThemesDir
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(fs) == 0 {
		return nil, errors.New("未找到任何主题文件")
	}

	for _, file := range fs {
		if file.IsDir() && (file.Name() == conf.Theme) {
			return loadTheme(path, conf)
		}
	}

	return nil, &helper.FieldError{
		Message: "不存在",
		Field:   "theme",
		File:    path.MetaConfigFile,
	}
}

// 加载主题
//
// id 主题当前目录名称
func loadTheme(path *path.Path, conf *config) (*Theme, error) {
	p := path.ThemeMetaPath(conf.Theme)

	theme := &Theme{}
	if err := helper.LoadYAMLFile(p, theme); err != nil {
		return nil, err
	}

	theme.ID = conf.Theme
	theme.longDateFormat = conf.LongDateFormat
	theme.shortDateFormat = conf.ShortDateFormat

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

// ExecuteTemplate 渲染指定的模块并输出到 w
func (d *Data) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	return d.Theme.template.ExecuteTemplate(w, name, data)
}

// 编译主题的模板。
func (d *Data) compileTemplate() error {
	snippets, err := d.snippetsTemplate()
	if err != nil {
		return err
	}

	// 编译模板
	d.Theme.template, err = snippets.Clone()
	if err != nil {
		return err
	}

	path := d.path.ThemesPath(d.Theme.ID, "*"+vars.TemplateExtension)
	_, err = d.Theme.template.ParseGlob(path)
	if err != nil {
		return err
	}

	// 检测模板名称是否在模板中真实存在
	// 模板定义未必是按文件分的，所以不能简单地判断文件是否存在
	templates := d.templatesName()
	for _, tpl := range templates {
		if nil == d.Theme.template.Lookup(tpl) {
			return fmt.Errorf("模板 %s 未定义", tpl)
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
		"ldate":    d.Theme.longDate,
		"sdate":    d.Theme.shortDate,
		"rfc3339":  rfc3339Date,
		"themeURL": func(p string) string { return vars.ThemeURL(p) },
	}

	return template.New("snippets").
		Funcs(funs).
		ParseGlob(d.path.ThemesPath("*" + vars.TemplateExtension))
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

func (theme *Theme) longDate(t time.Time) interface{} {
	return t.Format(theme.longDateFormat)
}

func (theme *Theme) shortDate(t time.Time) interface{} {
	return t.Format(theme.shortDateFormat)
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

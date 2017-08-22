// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"time"

	"github.com/caixw/typing/vars"
)

// 模板文件的扩展名
const templateExtension = ".html"

// 肯定存在的模板类型，检测模板是否存在时，会用到。
// 该值会被 checkTemplates 改变，不能用于其它地方。
var templates = []string{
	"posts",
	"post",
	"tags",
	"tag",
	"links",
	"archives",
	"search",
}

// 模板的扩展名，在主题目录下，以下扩展名的文件，不会被展示
var ignoreThemeFileExts = []string{
	templateExtension,
	".yaml",
	".yml",
}

func isIgnoreThemeFile(file string) bool {
	ext := filepath.Ext(file)

	for _, v := range ignoreThemeFileExts {
		if ext == v {
			return true
		}
	}

	return false
}

// 编译主题的模板。
func (client *Client) compileTemplate() error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"ldate":    client.longDateFormat,
		"sdate":    client.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return vars.ThemesURL(p) },
	}

	tpl, err := template.New("client").
		Funcs(funcMap).
		ParseGlob(filepath.Join(client.data.Theme.Path, "*"+templateExtension))
	if err != nil {
		return err
	}
	client.template = tpl

	return client.checkTemplatesExists()
}

// 检测模板名称是否在模板中真实存在
func (client *Client) checkTemplatesExists() error {
	for _, post := range client.data.Posts {
		for _, tpl := range templates {
			if tpl != post.Template {
				templates = append(templates, post.Template)
			}
		}
	}

	// 模板定义未必是按文件分的，所以不能简单地判断文件是否存在
	for _, tpl := range templates {
		if nil != client.template.Lookup(tpl) {
			continue
		}

		return fmt.Errorf("模板 %s 未定义", tpl)
	}

	return nil
}

func rfc3339DateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func (client *Client) longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(client.data.Config.LongDateFormat)
}

func (client *Client) shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(client.data.Config.ShortDateFormat)
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

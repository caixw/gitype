// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"html/template"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/caixw/typing/vars"
)

// 加载主题的模板。
// dir 模板所在的目录。
func (c *Client) initTemplate() error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"ldate":    c.longDateFormat,
		"sdate":    c.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return path.Join(vars.Themes, p) },
	}

	tpl, err := template.New("").
		Funcs(funcMap).
		ParseGlob(filepath.Join(c.path.ThemesDir, c.data.Config.Theme, "*.html"))
	if err != nil {
		return err
	}
	c.tpl = tpl
	return nil
}

func rfc3339DateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func (c *Client) longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(c.data.Config.LongDateFormat)
}

func (c *Client) shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(c.data.Config.ShortDateFormat)
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

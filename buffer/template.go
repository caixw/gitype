// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package buffer

import (
	"html/template"
	"path/filepath"
	"regexp"
	"time"

	"github.com/caixw/typing/vars"
)

// 编译主题的模板。
func (b *Buffer) compileTemplate() error {
	funcMap := template.FuncMap{
		"strip":    stripTags,
		"html":     htmlEscaped,
		"ldate":    b.longDateFormat,
		"sdate":    b.shortDateFormat,
		"rfc3339":  rfc3339DateFormat,
		"themeURL": func(p string) string { return vars.ThemesURL(p) },
	}

	tpl, err := template.New("").
		Funcs(funcMap).
		ParseGlob(filepath.Join(b.Data.Theme.Path, "*.html"))
	if err != nil {
		return err
	}
	b.Template = tpl

	return nil
}

func rfc3339DateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(time.RFC3339)
}

func (b *Buffer) longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(b.Data.Config.LongDateFormat)
}

func (b *Buffer) shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(b.Data.Config.ShortDateFormat)
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

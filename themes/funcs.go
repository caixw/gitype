// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/caixw/typing/core"
)

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}

func longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(opt.LongDateFormat)
}

func shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(opt.ShortDateFormat)
}

func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

func avatarImage(email string) interface{} {
	url := "http://www.gravatar.com/avatar/" + core.MD5(strings.ToLower(email))
	// TODO 将选项添加到options中
	return url + "?s=96&d=mm&r=g"
}

var funcMap = template.FuncMap{
	"html":   htmlEscaped,
	"ldate":  longDateFormat,
	"sdate":  shortDateFormat,
	"strip":  stripTags,
	"avatar": avatarImage,
}

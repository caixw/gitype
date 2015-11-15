// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"errors"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/caixw/typing/core"
)

// 去掉所有的标签信息
var stripExpr = regexp.MustCompile("</?[^</>]+/?>")

// 过滤标签。
func stripTags(html string) string {
	return stripExpr.ReplaceAllString(html, "")
}

// 根据options中的格式显示长日期
func longDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(opt.LongDateFormat)
}

// 根据options中的格式显示短日期
func shortDateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(opt.ShortDateFormat)
}

// 将内容显示为html内容
func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

// 转换成gravatar头像
func avatarImage(email string) interface{} {
	url := "https://secure.gravatar.com/avatar/" + core.MD5(strings.ToLower(email))
	// TODO 将选项添加到options中
	return url + "?s=96&d=mm&r=g"
}

// 将给定的文件转换成相对主题目录的URL
func themeURL(url string) (interface{}, error) {
	if strings.Index(url, "..") >= 0 { // 不能只判断单个点号，文件可能带后缀名
		return nil, errors.New("路径中不能带点号")
	}

	return cfg.ThemeURLPrefix + "/" + url, nil
}

var funcMap = template.FuncMap{
	"html":     htmlEscaped,
	"ldate":    longDateFormat,
	"sdate":    shortDateFormat,
	"strip":    stripTags,
	"avatar":   avatarImage,
	"themeURL": themeURL,
}

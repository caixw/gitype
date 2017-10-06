// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/vars"
)

// 一些默认值的定义
// 可以使用 %s 和 %t 两个变量
const (
	tagTitle      = "标签：%s"
	tagsTitle     = "标签"
	archivesTitle = "归档"
	searchTitle   = "搜索：%s"
	linksTitle    = "友情链接"
	postTitle     = "%s"
	//postsTitle    = "第 %d 页"
	//indexTitle    = ""
)

// Page 页面的自定义内容
type Page struct {
	Title       string `yaml:"title"`
	Keywords    string `yaml:"keywords"`
	Description string `yaml:"description"`
}

func (conf *config) initPages() *helper.FieldError {
	if conf.Pages == nil {
		conf.Pages = make(map[string]*Page, 10)
	}

	if conf.Pages[vars.PageArchives] == nil {
		conf.Pages[vars.PageArchives] = &Page{
			Title: archivesTitle,
		}
	}
	conf.Pages[vars.PageArchives].Title += conf.Title

	return nil
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"github.com/caixw/gitype/vars"
)

// 一些默认值的定义
// 可以使用 %s 和 %t 两个变量
const (
	tagTitle      = "标签：%s | %t"
	tagsTitle     = "标签 | %t"
	archivesTitle = "归档 | %t"
	searchTitle   = "搜索：%s | %t"
	linksTitle    = "友情链接 | %t"
	postTitle     = "%s | %t"
	homeTitle     = "%t"
)

// Page 页面的自定义内容
type Page struct {
	Title       string `yaml:"title"`
	Keywords    string `yaml:"keywords"`
	Description string `yaml:"description"`
}

func (conf *config) initPages() {
	if conf.Pages == nil {
		conf.Pages = make(map[string]*Page, 10)
	}

	if conf.Pages[vars.PageTag] == nil {
		conf.Pages[vars.PageTag] = &Page{
			Title: tagTitle,
		}
	}

	if conf.Pages[vars.PageTags] == nil {
		conf.Pages[vars.PageTags] = &Page{
			Title: tagsTitle,
		}
	}

	if conf.Pages[vars.PageArchives] == nil {
		conf.Pages[vars.PageArchives] = &Page{
			Title: archivesTitle,
		}
	}

	if conf.Pages[vars.PageSearch] == nil {
		conf.Pages[vars.PageSearch] = &Page{
			Title: searchTitle,
		}
	}

	if conf.Pages[vars.PageLinks] == nil {
		conf.Pages[vars.PageLinks] = &Page{
			Title: linksTitle,
		}
	}

	if conf.Pages[vars.PagePost] == nil {
		conf.Pages[vars.PagePost] = &Page{
			Title: postTitle,
		}
	}

	if conf.Pages[vars.PagePost] == nil {
		conf.Pages[vars.PagePost] = &Page{
			Title: postTitle,
		}
	}

	if conf.Pages[vars.PagePosts] == nil {
		conf.Pages[vars.PagePosts] = &Page{
			Title: homeTitle,
		}
	}
	if conf.Pages[vars.PageIndex] == nil {
		conf.Pages[vars.PageIndex] = conf.Pages[vars.PagePosts]
	}
}

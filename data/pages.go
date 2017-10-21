// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strings"

	"github.com/caixw/gitype/vars"
)

// 一些默认值的定义
const (
	tagTitle      = "标签：" + vars.ContentPlaceholder + " | " + vars.TitlePlaceholder
	tagsTitle     = "标签 | " + vars.TitlePlaceholder
	archivesTitle = "归档 | " + vars.TitlePlaceholder
	searchTitle   = "搜索：" + vars.ContentPlaceholder + " | " + vars.TitlePlaceholder
	linksTitle    = "友情链接 | " + vars.TitlePlaceholder
	postTitle     = vars.ContentPlaceholder + " | " + vars.TitlePlaceholder
	homeTitle     = vars.TitlePlaceholder
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

	for _, page := range conf.Pages {
		page.Title = conf.replaceTitle(page.Title)
		page.Keywords = conf.replaceTitle(page.Keywords)
		page.Description = conf.replaceTitle(page.Description)
	}

	if conf.Pages[vars.PageIndex] == nil {
		conf.Pages[vars.PageIndex] = conf.Pages[vars.PagePosts]
	}
}

// 替换标题中的 %title% 内容为 conf.Title
func (conf *config) replaceTitle(title string) string {
	if strings.Index(title, vars.TitlePlaceholder) < 0 {
		return title
	}

	return strings.Replace(title, vars.TitlePlaceholder, conf.Title, -1)
}

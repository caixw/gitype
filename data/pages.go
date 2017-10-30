// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import "github.com/caixw/gitype/vars"

// 默认标题的定义
const (
	tagTitle      = "标签：" + vars.ContentPlaceholder
	tagsTitle     = "标签"
	archivesTitle = "归档"
	searchTitle   = "搜索：" + vars.ContentPlaceholder
	linksTitle    = "友情链接"
	postTitle     = vars.ContentPlaceholder
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
	ps := conf.Pages

	if ps[vars.PageTag] == nil {
		ps[vars.PageTag] = &Page{}
	}
	if ps[vars.PageTags] == nil {
		ps[vars.PageTags] = &Page{}
	}
	if ps[vars.PageArchives] == nil {
		ps[vars.PageArchives] = &Page{}
	}
	if ps[vars.PageSearch] == nil {
		ps[vars.PageSearch] = &Page{}
	}
	if ps[vars.PageLinks] == nil {
		ps[vars.PageLinks] = &Page{}
	}
	if ps[vars.PagePost] == nil {
		ps[vars.PagePost] = &Page{}
	}
	if ps[vars.PagePost] == nil {
		ps[vars.PagePost] = &Page{}
	}
	if ps[vars.PagePosts] == nil {
		ps[vars.PagePosts] = &Page{}
	}

	if len(ps[vars.PageTag].Title) == 0 {
		ps[vars.PageTag].Title = tagTitle
	}

	if len(ps[vars.PageTags].Title) == 0 {
		ps[vars.PageTags].Title = tagsTitle
	}

	if len(ps[vars.PageArchives].Title) == 0 {
		ps[vars.PageArchives].Title = archivesTitle
	}

	if len(ps[vars.PageSearch].Title) == 0 {
		ps[vars.PageSearch].Title = searchTitle
	}

	if len(ps[vars.PageLinks].Title) == 0 {
		ps[vars.PageLinks].Title = linksTitle
	}

	if len(ps[vars.PagePost].Title) == 0 {
		ps[vars.PagePost].Title = postTitle
	}

	suffix := conf.TitleSeparator + conf.Title
	for _, page := range conf.Pages {
		if len(page.Title) == 0 { // 没有内容，则直接使用网站标
			page.Title = conf.Title
		} else {
			page.Title += suffix
		}
	}

	ps[vars.PageIndex] = conf.Pages[vars.PagePosts]
}

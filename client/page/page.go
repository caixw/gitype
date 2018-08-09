// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package page 定义 html 页面需要的内容
package page

import "github.com/caixw/gitype/data"

// Page 用于描述一个页面的所有无素
type Page struct {
	Site *Site

	Title       string       // 文章标题
	Subtitle    string       // 副标题
	Canonical   string       // 当前页的唯一链接
	Keywords    string       // meta.keywords 的值
	Description string       // meta.description 的值
	PrevPage    *data.Link   // 前一页
	NextPage    *data.Link   // 下一页
	Type        string       // 当前页面类型
	Charset     string       // 当前页的字符集
	Author      *data.Author // 作者
	License     *data.Link   // 当前页的版本信息，可以为空

	// 以下内容，仅在对应的页面才会有内容
	Q        string          // 搜索关键字
	Tag      *data.Tag       // 标签详细页面，非标签详细页，则为空
	Posts    []*data.Post    // 文章列表，仅标签详情页和搜索页用到。
	Post     *data.Post      // 文章详细内容，仅文章页面用到。
	Archives []*data.Archive // 归档
}

// Page 生成 Page 实例
func (site *Site) Page(typ string, d *data.Data) *Page {
	return &Page{
		Site: site,

		Subtitle: d.Subtitle,
		Type:     typ,
		Author:   d.Author,
		License:  d.License,
	}
}

// Next 产生 NextPage 内容
func (p *Page) Next(url, text string) {
	p.NextPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "next",
	}
}

// Prev 产生 PrevPage 内容
func (p *Page) Prev(url, text string) {
	p.PrevPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "prev",
	}
}

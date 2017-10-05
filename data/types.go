// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/issue9/is"
)

// Feed RSS、Atom、Sitemap 和 Opensearch 的配置内容
type Feed struct {
	Title   string // 标题，一般出现在 html>head>link.title 属性中
	URL     string // 地址，不能包含域名
	Type    string // mime type
	Content []byte // 实际的内容
}

// Author 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

// Link 描述链接的内容
type Link struct {
	// 链接对应的图标。可以是字体图标或是图片链接，模板根据情况自动选择。
	Icon  string `yaml:"icon,omitempty"`
	Title string `yaml:"title,omitempty"` // 链接的 title 属性
	Rel   string `yaml:"rel,omitempty"`   // 链接的 rel 属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
	Type  string `yaml:"type,omitempty"`  // 链接的类型，一般用于 a 和 link 标签的 type 属性
}

// Icon 表示网站图标，比如 html>head>link.rel="short icon"
type Icon struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type"` // mime type
	Sizes string `yaml:"sizes"`
}

func loadLinks(path *path.Path) ([]*Link, error) {
	links := make([]*Link, 0, 20)
	if err := helper.LoadYAMLFile(path.MetaLinksFile, &links); err != nil {
		return nil, err
	}

	for index, link := range links {
		if err := link.sanitize(); err != nil {
			err.File = path.MetaLinksFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return nil, err
		}
	}

	return links, nil
}

func (icon *Icon) sanitize() *helper.FieldError {
	if len(icon.URL) == 0 {
		return &helper.FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (link *Link) sanitize() *helper.FieldError {
	if len(link.Text) == 0 {
		return &helper.FieldError{Field: "text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &helper.FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (author *Author) sanitize() *helper.FieldError {
	if len(author.Name) == 0 {
		return &helper.FieldError{Field: "name", Message: "不能为空"}
	}

	if len(author.URL) > 0 && !is.URL(author.URL) {
		return &helper.FieldError{Field: "url", Message: "不是一个正确的 URL"}
	}

	if len(author.Avatar) > 0 && !is.URL(author.Avatar) {
		return &helper.FieldError{Field: "avatar", Message: "不是一个正确的 URL"}
	}

	if len(author.Email) > 0 && !is.Email(author.Email) {
		return &helper.FieldError{Field: "email", Message: "不是一个正确的 Email"}
	}

	return nil
}

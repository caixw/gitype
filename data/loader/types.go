// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"strconv"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/issue9/is"
)

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

// Theme 表示主题信息
type Theme struct {
	ID          string  `yaml:"-"`    // 唯一 ID，即当前目录名称
	Name        string  `yaml:"name"` // 名称，不必唯一，可以与 ID 值不同。
	Version     string  `yaml:"version"`
	Description string  `yaml:"description"`
	URL         string  `yaml:"url,omitempty"`
	Author      *Author `yaml:"author"`

	// 需要被 service worker 缓存的内容。
	// 如果是带 https 开头的 URL，则直接使用，
	// 如果是不以 https 开头的 URL，则会被映射到当前主题下。
	Assets []string `yaml:"assets,omitempty"`
}

// LoadLinks 加载友情链接的内容
func LoadLinks(path *path.Path) ([]*Link, error) {
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

// LoadTheme 加载当前的主题
//
// name 为主题的唯一名称
func LoadTheme(path *path.Path, name string) (*Theme, error) {
	p := path.ThemeMetaPath(name)

	theme := &Theme{}
	if err := helper.LoadYAMLFile(p, theme); err != nil {
		return nil, err
	}
	theme.ID = name

	if len(theme.Name) == 0 {
		return nil, &helper.FieldError{File: path.ThemeMetaPath(theme.ID), Message: "不能为空", Field: "name"}
	}

	if theme.Author != nil {
		if err := theme.Author.sanitize(); err != nil {
			err.Field = path.ThemeMetaPath(theme.ID)
			return nil, err
		}
	}

	return theme, nil
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

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"strconv"

	"github.com/caixw/typing/vars"
)

// Feed RSS、Atom 和 Opensearch 等的配置内容
type Feed struct {
	Title   string // 标题，一般出现在 html>head>link.title 属性中
	URL     string // 地址，不能包含域名
	Type    string // mimeType
	Content []byte // 的实际内容
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
	// 链接对应的图标名称，fontawesome 图标名称，不用带 fa- 前缀。
	// 也有可能是链接，模板根据情况自动选择。
	Icon  string `yaml:"icon,omitempty"`
	Title string `yaml:"title,omitempty"` // 链接的 title 属性
	Rel   string `yaml:"rel,omitempty"`   // 链接的 rel 属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
}

// Icon 表示程序图标
type Icon struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type"` // mime type
	Sizes string `yaml:"sizes"`
}

// FieldError 表示加载文件出错时，具体的错误信息
type FieldError struct {
	File    string // 所在文件
	Message string // 错误信息
	Field   string // 所在的字段
}

func loadLinks(path *vars.Path) ([]*Link, error) {
	links := make([]*Link, 0, 20)
	if err := loadYamlFile(path.MetaLinksFile, &links); err != nil {
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

func (err *FieldError) Error() string {
	return fmt.Sprintf("在文件 %s 中的 %s 字段发生错误：%s", err.File, err.Field, err.Message)
}

func (icon *Icon) sanitize() *FieldError {
	if len(icon.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (link *Link) sanitize() *FieldError {
	if len(link.Text) == 0 {
		return &FieldError{Field: "text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (author *Author) sanitize() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "name", Message: "不能为空"}
	}

	return nil
}

func (o *outdatedConfig) sanitize() *FieldError {
	if o.Type != outdatedTypeCreated && o.Type != outdatedTypeModified {
		return &FieldError{Message: "无效的值", Field: "outdated.type"}
	}

	if len(o.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "outdated.content"}
	}

	return nil
}

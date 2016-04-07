// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Link 用于描述链接的内容
type Link struct {
	Icon  string `yaml:"icon,omitempty"`  // 链接对应的图标名称，fontawesome图标名称，不用带fa-前缀。
	Title string `yaml:"title,omitempty"` // 链接的title属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
}

func (d *Data) loadLinks(p string) error {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	links := make([]*Link, 0, 20)
	if err = yaml.Unmarshal(data, &links); err != nil {
		return err
	}

	// 检测错误
	for index, link := range links {
		if err := link.check(); err != nil {
			err.File = "links.yaml"
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	d.Links = links
	return nil
}

func (link *Link) check() *MetaError {
	if len(link.Text) == 0 {
		return &MetaError{Field: "Text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &MetaError{Field: "URL", Message: "不能为空"}
	}

	return nil
}

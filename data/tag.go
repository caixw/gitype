// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// 描述标签信息
type Tag struct {
	Slug      string `yaml:"slug"`
	Title     string `yaml:"title"`
	Color     string `yaml:"color,omitempty"` // 未指定，则继承父容器
	Content   string `yaml:"content"`
	Count     int    `yaml:"-"` // 文章计数
	Premalink string `yaml:"-"`
}

func (d *Data) loadTags(p string) error {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	tags := make([]*Tag, 0, 100)
	if err = yaml.Unmarshal(data, &tags); err != nil {
		return err
	}
	for index, tag := range tags {
		if len(tag.Slug) == 0 {
			return fmt.Errorf("第[%v]个标签未指slug", index)
		}

		if len(tag.Title) == 0 {
			return fmt.Errorf("第[%v]个标签未指title", index)
		}

		if len(tag.Content) == 0 {
			return fmt.Errorf("第[%v]个标签未指content", index)
		}

		tag.Premalink = path.Join(d.URLS.Root, d.URLS.Tag, tag.Slug+d.URLS.Suffix)
	}
	d.Tags = tags
	return nil
}

// 查找指定名称的标签。
// 若返回nil，则表示该标签不存在。
func (d *Data) FindTag(slug string) *Tag {
	for _, tag := range d.Tags {
		if tag.Slug == slug {
			return tag
		}
	}

	return nil
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// 描述标签信息
type Tag struct {
	Slug      string `yaml:"slug"`
	Title     string `yaml:"title"`
	Color     string `yaml:"color,omitempty"`
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
	for _, tag := range tags {
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

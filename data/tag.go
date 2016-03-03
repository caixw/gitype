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
	Slug      string  `yaml:"slug"`            //  唯一名称
	Title     string  `yaml:"title"`           // 名称
	Color     string  `yaml:"color,omitempty"` // 标签颜色。若未指定，则继承父容器
	Content   string  `yaml:"content"`         // 对该标签的详细描述
	Posts     []*Post `yaml:"-"`               // 关联的文章
	Permalink string  `yaml:"-"`               // 唯一链接
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

		tag.Posts = make([]*Post, 0, 10)
		tag.Permalink = path.Join(d.URLS.Root, d.URLS.Tag, tag.Slug+d.URLS.Suffix)
	}
	d.Tags = tags
	return nil
}

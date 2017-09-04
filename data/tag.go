// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"strconv"
	"time"

	"github.com/caixw/typing/vars"
)

// Tag 描述标签信息
type Tag struct {
	Slug        string    `yaml:"slug"`            // 唯一名称
	Title       string    `yaml:"title"`           // 名称
	Color       string    `yaml:"color,omitempty"` // 标签颜色。若未指定，则继承父容器
	Content     string    `yaml:"content"`         // 对该标签的详细描述
	Posts       []*Post   `yaml:"-"`               // 关联的文章
	Keywords    string    `yaml:"-"`               // meta.keywords 标签的内容，如果为空，使用 Title 属性的值
	Description string    `yaml:"-"`               // meta.description 标签的内容，若为空，则为 Config.Description
	Modified    time.Time `yaml:"-"`               // 所有文章中最迟修改的
	Permalink   string    `yaml:"-"`               // 唯一链接，指向第一页
}

func loadTags(path *vars.Path) ([]*Tag, error) {
	tags := make([]*Tag, 0, 100)
	if err := loadYamlFile(path.MetaTagsFile, &tags); err != nil {
		return nil, err
	}

	for index, tag := range tags {
		if err := tag.sanitize(); err != nil {
			err.File = path.MetaTagsFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return nil, err
		}
	}

	// 等待其它检测完成，再检查是否存在同名的
	if err := checkTagsDup(tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// 检测是否存在同名的标签
func checkTagsDup(tags []*Tag) error {
	count := func(slug string) (cnt int) {
		for _, tag := range tags {
			if tag.Slug == slug {
				cnt++
			}
		}
		return cnt
	}

	for _, tag := range tags {
		if count(tag.Slug) > 1 {
			return errors.New("存在同名的标签：" + tag.Slug)
		}
	}

	return nil
}

func (tag *Tag) sanitize() *FieldError {
	if len(tag.Slug) == 0 {
		return &FieldError{Message: "不能为空", Field: "slug"}
	}

	if len(tag.Title) == 0 {
		return &FieldError{Message: "不能为空", Field: "title"}
	}

	if len(tag.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "content"}
	}

	tag.Posts = make([]*Post, 0, 100)

	tag.Permalink = vars.TagURL(tag.Slug, 1)

	tag.Keywords = tag.Title
	if tag.Title != tag.Slug {
		tag.Keywords += ","
		tag.Keywords += tag.Slug
	}

	tag.Description = "标签" + tag.Title + "的介绍"

	return nil
}

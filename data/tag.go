// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strings"
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

// Tag 描述标签信息
//
// 标签系统同时包含了标签和专题两个方面，默认情况下为标签，
// 当将 Series 指定为 true 时，表示这是一个专题。
type Tag struct {
	loader.Tag

	HTMLTitle string    `yaml:"-"` // 用于网页的标题
	Posts     []*Post   `yaml:"-"` // 关联的文章
	Keywords  string    `yaml:"-"` // meta.keywords 标签的内容，如果为空，使用 Tag.Title 属性的值
	Modified  time.Time `yaml:"-"` // 所有文章中最迟修改的
	Permalink string    `yaml:"-"` // 唯一链接，指向第一页

	// 用于搜索的副本内容，会全部转换成小写
	SearchTitle string
}

func loadTags(path *path.Path, conf *loader.Config) ([]*Tag, error) {
	tags, err := loader.LoadTags(path)
	if err != nil {
		return nil, err
	}

	ret := make([]*Tag, 0, len(tags))
	p := conf.Pages[vars.PageTag]
	for _, tag := range tags {
		keywords := tag.Title
		if tag.Title != tag.Slug {
			keywords = keywords + "," + tag.Slug
		}

		ret = append(ret, &Tag{
			Tag:         *tag,
			Posts:       make([]*Post, 0, 100),
			Permalink:   vars.TagURL(tag.Slug, 1),
			SearchTitle: strings.ToLower(tag.Title),
			Keywords:    keywords,
			Modified:    conf.Uptime,
			HTMLTitle:   helper.ReplaceContent(p.Title, tag.Title),
		})
	}

	return ret, nil
}

// 分离标签和专题的列表
func splitTags(tags []*Tag) (ts []*Tag, series []*Tag) {
	ts = make([]*Tag, 0, len(tags))
	series = make([]*Tag, 0, len(tags))

	for _, tag := range tags {
		if tag.Series {
			series = append(series, tag)
		} else {
			ts = append(ts, tag)
		}
	}

	return ts, series
}

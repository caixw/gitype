// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"sort"
	"time"

	"github.com/caixw/gitype/helper"
)

// 归档的类型
const (
	archiveTypeYear  = "year"
	archiveTypeMonth = "month"
)

// 归档的排序方式
const (
	archiveOrderDesc = "desc"
	archiveOrderAsc  = "asc"
)

// Archive 表示某一时间段的存档信息
type Archive struct {
	date  time.Time // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string    // 当前存档的标题
	Posts []*Post   // 当前存档的文章列表
}

// 存档页的配置内容
type archiveConfig struct {
	Order  string `yaml:"order"`            // 排序方式
	Type   string `yaml:"type,omitempty"`   // 存档的分类方式，可以按年或是按月
	Format string `yaml:"format,omitempty"` // 标题的格式化字符串
}

func (d *Data) buildArchives(conf *config) error {
	archives := make([]*Archive, 0, 10)

	for _, post := range d.Posts {
		t := post.Created
		var date time.Time

		switch conf.Archive.Type {
		case archiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location())
		case archiveTypeYear:
			date = time.Date(t.Year(), 2, 0, 0, 0, 0, 0, t.Location())
		default:
			return &helper.FieldError{File: d.path.MetaConfigFile, Field: "archive.type", Message: "无效的取值"}
		}

		found := false
		for _, archive := range archives {
			if archive.date.Equal(date) {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			archives = append(archives, &Archive{
				date:  date,
				Title: date.Format(conf.Archive.Format),
				Posts: []*Post{post},
			})
		}
	} // end for

	sort.SliceStable(archives, func(i, j int) bool {
		if conf.Archive.Order == archiveOrderDesc {
			return archives[i].date.After(archives[j].date)
		}
		return archives[i].date.Before(archives[j].date)
	})

	d.Archives = archives

	return nil
}

func (a *archiveConfig) sanitize() *helper.FieldError {
	if len(a.Type) == 0 {
		a.Type = archiveTypeYear
	} else {
		if a.Type != archiveTypeMonth && a.Type != archiveTypeYear {
			return &helper.FieldError{Message: "取值不正确", Field: "archive.type"}
		}
	}

	if len(a.Order) == 0 {
		a.Order = archiveOrderDesc
	} else {
		if a.Order != archiveOrderAsc && a.Order != archiveOrderDesc {
			return &helper.FieldError{Message: "取值不正确", Field: "archive.order"}
		}
	}

	return nil
}

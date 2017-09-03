// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"sort"
	"time"
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

// Archive 存档信息
type Archive struct {
	date  int64
	Title string
	Posts []*Post
}

// archiveConfig 存档页的配置内容
type archiveConfig struct {
	Order  string `yaml:"order"`            // 排序方式
	Type   string `yaml:"type,omitempty"`   // 存档的分类方式，可以为 year 或是 month
	Format string `yaml:"format,omitempty"` // 标题的格式化字符串
}

func (d *Data) buildArchives() error {
	archives := make([]*Archive, 0, 10)

	for _, post := range d.Posts {
		t := post.Created
		var date int64

		switch d.Config.Archive.Type {
		case archiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location()).Unix()
		case archiveTypeYear:
			date = time.Date(t.Year(), 2, 0, 0, 0, 0, 0, t.Location()).Unix()
		}

		found := false
		for _, archive := range archives {
			if archive.date == date {
				archive.Posts = append(archive.Posts, post)
				found = true
				break
			}
		}
		if !found {
			archives = append(archives, &Archive{
				date:  date,
				Title: time.Unix(date, 0).Format(d.Config.Archive.Format),
				Posts: []*Post{post},
			})
		}
	} // end for

	sort.SliceStable(archives, func(i, j int) bool {
		less := archives[i].date > archives[j].date
		if d.Config.Archive.Order == archiveOrderDesc {
			less = !less
		}

		return less
	})

	d.Archives = archives

	return nil
}

func (a *archiveConfig) sanitize() *FieldError {
	if len(a.Type) == 0 {
		a.Type = archiveTypeYear
	} else {
		if a.Type != archiveTypeMonth && a.Type != archiveTypeYear {
			return &FieldError{Message: "取值不正确", Field: "archive.type"}
		}
	}

	if len(a.Order) == 0 {
		a.Order = archiveOrderAsc
	} else {
		if a.Order != archiveOrderAsc && a.Order != archiveOrderDesc {
			return &FieldError{Message: "取值不正确", Field: "archive.order"}
		}
	}

	return nil
}

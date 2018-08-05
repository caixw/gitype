// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"sort"
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/helper"
)

// Archive 表示某一时间段的存档信息
type Archive struct {
	date  time.Time // 当前存档的一个日期值，可用于生成 Title 和排序用，具体取值方式，可自定义
	Title string    // 当前存档的标题
	Posts []*Post   // 当前存档的文章列表
}

func (d *Data) buildArchives(conf *loader.Config) error {
	archives := make([]*Archive, 0, 10)

	for _, post := range d.Posts {
		t := post.Created
		var date time.Time

		switch conf.Archive.Type {
		case loader.ArchiveTypeMonth:
			date = time.Date(t.Year(), t.Month(), 2, 0, 0, 0, 0, t.Location())
		case loader.ArchiveTypeYear:
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
		if conf.Archive.Order == loader.ArchiveOrderDesc {
			return archives[i].date.After(archives[j].date)
		}
		return archives[i].date.Before(archives[j].date)
	})

	d.Archives = archives

	return nil
}

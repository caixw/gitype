// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/caixw/typing/util"
	"github.com/issue9/conv"
	"github.com/issue9/orm"
)

// 向数据库写入初始内容。
func Install(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// option
	return fillOptions(db)
}

// 将Options中的每字段转换成一个map结构，方便其它工具将其转换成sql内容。
//  options.PageSize=5 ==> {"group":"system", "key":"pageSize", "value":"5"}
func (opt *Options) toMaps() ([]map[string]string, error) {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	l := t.NumField()
	maps := make([]map[string]string, 0, l)

	for i := 0; i < l; i++ {
		tags := strings.Split(t.Field(i).Tag.Get("options"), ",")
		if len(tags) != 2 {
			return nil, fmt.Errorf("len(tags)!=2 @ %v", t.Field(i).Name)
		}

		val, err := conv.String(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		maps = append(maps, map[string]string{
			"group": tags[0],
			"key":   tags[1],
			"value": val,
		})
	}

	return maps, nil
}

// 将一个默认的options值填充到数据库中。
func fillOptions(db *orm.DB) error {
	now := time.Now().Unix()
	opt := &Options{
		SiteURL:     "http://localhost:8080/",
		SiteName:    "typing blog",
		SecondTitle: "副标题",
		ScreenName:  "typing",
		Password:    util.HashPassword("123"),
		Keywords:    "typing",
		Description: "typing-极简的博客系统",
		Suffix:      ".html",

		Uptime:               now,
		LastUpdated:          now,
		CommentsSize:         0,
		WattingCommentsSize:  0,
		ApprovedCommentsSize: 0,
		SpamCommentsSize:     0,
		PostsSize:            0,
		PublishedPostsSize:   0,
		DraftPostsSize:       0,
		LastLogin:            0,
		LastIP:               "",
		LastAgent:            "",

		PageSize:        20,
		LongDateFormat:  "2006-01-02 15:04:05",
		ShortDateFormat: "2006-01-02",
		SidebarSize:     10,
		CommentOrder:    CommentOrderDesc,

		PostsChangefreq: "never",
		TagsChangefreq:  "daily",
		PostsPriority:   0.9,
		TagsPriority:    0.4,
		RssSize:         20,

		Theme: "default",
	}

	maps, err := opt.toMaps()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sql := "INSERT INTO #options ({key},{group},{value}) VALUES(?,?,?)"
	stmt, err := tx.Prepare(true, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range maps {
		_, err := stmt.Exec(item["key"], item["group"], item["value"])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/issue9/conv"
	"github.com/issue9/orm"
)

// 向数据库写入初始内容。
func Install(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// 创建表
	if err := createTables(db); err != nil {
		return err
	}

	// option
	if err := fillOptions(db); err != nil {
		return err
	}

	// tags
	if err := fillTags(db); err != nil {
		return err
	}

	// post
	now := time.Now().Unix()
	post := &Post{
		Title:    "第一篇日志",
		Content:  "<p>这是你的第一篇日志</p>",
		State:    PostStatePublished,
		Created:  now,
		Modified: now,
	}
	if _, err := db.Insert(post); err != nil {
		return err
	}

	// comment
	comment := &Comment{
		PostID:     1,
		Content:    "<p>沙发</p>",
		AuthorName: "游客",
		State:      CommentStateWaiting,
	}
	if _, err := db.Insert(comment); err != nil {
		return err
	}

	// relationship
	if _, err := db.Insert(&Relationship{TagID: 1, PostID: 1}); err != nil {
		return err
	}
	if _, err := db.Insert(&Relationship{TagID: 2, PostID: 1}); err != nil {
		return err
	}

	return nil
}

// 创建所有表结构
func createTables(db *orm.DB) error {
	// 创建表
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = tx.MultCreate(
		&Option{},
		&Comment{},
		&Tag{},
		&Post{},
		&Relationship{},
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 填充标签
func fillTags(db *orm.DB) error {
	tags := []*Tag{
		{Name: "default", Title: "默认标签", Description: "这是系统产生的默认标签"},
		{Name: "tag1", Title: "标签一", Description: "tag1"},
		{Name: "tag2", Title: "标签二", Description: "tag2"},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := tx.InsertMany(tags); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 将Options中的每字段转换成一个map结构，方便其它工具将其转换成sql内容。
//  options.PageSize=5 ==> {"group":"system", "key":"pageSize", "value":"5"}
func toMaps(opt *Options) ([]map[string]string, error) {
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
	opt := &Options{
		SiteURL:     "http://localhost:8080/",
		SiteName:    "typing blog",
		SecondTitle: "副标题",
		ScreenName:  "typing",
		Password:    HashPassword("123"),
		Keywords:    "typing",
		Description: "typing-极简的博客系统",
		Suffix:      ".html",
		Uptime:      time.Now().Unix(),

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

	maps, err := toMaps(opt)
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

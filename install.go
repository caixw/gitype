// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"errors"

	"github.com/caixw/typing/core"
	"github.com/issue9/orm"
)

// 向数据库写入初始内容。
func fillDB(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// option
	if err := db.Create(&option{}); err != nil {
		return err
	}

	if err := fillOptions(db); err != nil {
		return err
	}

	// meta
	if err := db.Create(&meta{}); err != nil {
		return err
	}
	metas := []*meta{
		{Name: "default", Title: "默认分类", Type: metaTypeCat, Order: 10, Parent: metaNoParent, Description: "所有添加的文章，默认添加此分类下。"},

		{Name: "tag1", Title: "标签一", Type: metaTypeTag, Description: "<h5>tag1</h5>"},
		{Name: "tag2", Title: "标签二", Type: metaTypeTag, Description: "<h5>tag2</h5>"},
		{Name: "tag3", Title: "标签三", Type: metaTypeTag, Description: "<h5>tag3</h5>"},
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if err := tx.InsertMany(metas); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	// post
	if err := db.Create(&post{}); err != nil {
		return err
	}

	if _, err := db.Insert(&post{Title: "第一篇日志", Content: "<p>这是你的第一篇日志</p>"}); err != nil {
		return err
	}

	// comment
	if err := db.Create(&comment{}); err != nil {
		return err
	}

	if _, err := db.Insert(&comment{PostID: 1, Content: "<p>沙发</p>", AuthorName: "游客"}); err != nil {
		return err
	}

	// relationship
	if err := db.Create(&relationship{}); err != nil {
		return err
	}
	if _, err := db.Insert(&relationship{MetaID: 1, PostID: 1}); err != nil {
		return err
	}

	return nil
}

func fillOptions(db *orm.DB) error {
	opt := &options{
		PageSize:   20,
		SiteName:   "typing blog",
		ScreenName: "typing",
		Password:   core.HashPassword(defaultPassword),
		Theme:      "default",
		Keywords:   "typing",
	}

	maps, err := opt.toMaps()
	if err != nil {
		return err
	}

	sql := "INSERT INTO #options ({key},{group},{value}) VALUES(?,?,?)"
	for _, item := range maps {
		_, err := db.Exec(true, sql, item["key"], item["group"], item["value"])
		if err != nil {
			return err
		}
	}
	return nil
}

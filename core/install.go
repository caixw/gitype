// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"time"

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

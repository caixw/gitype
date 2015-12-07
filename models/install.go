// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// models 定义了所有的数据模型。
package models

import (
	"errors"
	"time"

	"github.com/issue9/orm"
)

// 安装数据库和初始化默认的初始数据。
// options表的数据在options包中安装。
func Install(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// 创建表
	if err := createTables(db); err != nil {
		return err
	}

	// tags
	tag := &Tag{
		Name:        "default",
		Title:       "默认标签",
		Description: "这是系统产生的默认标签",
	}
	if _, err := db.Insert(tag); err != nil {
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

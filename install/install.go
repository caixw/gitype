// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 用于执行typing的安装。包含执行以下几个步骤的内容：
// 输出默认的配置文件；
// 输出默认的日志配置文件；
// 填充默认的数据到数据库；
package install

import (
	"errors"
	"time"

	"github.com/caixw/typing/models"
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
	post := &models.Post{
		Title:    "第一篇日志",
		Content:  "<p>这是你的第一篇日志</p>",
		State:    models.PostStatePublished,
		Created:  now,
		Modified: now,
	}
	if _, err := db.Insert(post); err != nil {
		return err
	}

	// comment
	comment := &models.Comment{
		PostID:     1,
		Content:    "<p>沙发</p>",
		AuthorName: "游客",
		State:      models.CommentStateWaiting,
	}
	if _, err := db.Insert(comment); err != nil {
		return err
	}

	// relationship
	if _, err := db.Insert(&models.Relationship{TagID: 1, PostID: 1}); err != nil {
		return err
	}
	if _, err := db.Insert(&models.Relationship{TagID: 2, PostID: 1}); err != nil {
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
		&models.Option{},
		&models.Comment{},
		&models.Tag{},
		&models.Post{},
		&models.Relationship{},
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
	tags := []*models.Tag{
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

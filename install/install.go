// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 用于执行typing的安装。包含执行以下几个步骤的内容：
// 输出默认的配置文件；
// 输出默认的日志配置文件；
// 填充默认的数据到数据库；
package install

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/install/static"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

// Install 执行安装过程。
func Install() bool {
	action := flag.String("init", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	flag.Parse()
	switch *action {
	case "config":
		if err := outputConfigFile(core.LogConfigPath, core.ConfigPath); err != nil {
			panic(err)
		}
		return true
	case "db":
		cfg, err := core.LoadConfig(core.ConfigPath)
		if err != nil {
			panic(err)
		}

		db, err := core.InitDB(cfg)
		defer db.Close()
		if err != nil {
			panic(err)
		}
		if err := fillDB(db); err != nil {
			panic(err)
		}
		return true
	} // end switch

	return false
}

// 用于输出配置文件到指定的位置。
// 目前包含了日志配置文件和程序本身的配置文件。
func outputConfigFile(logsConfigPath, configPath string) error {
	if err := ioutil.WriteFile(logsConfigPath, static.LogConfig, os.ModePerm); err != nil {
		return err
	}

	cfg := &core.Config{
		Core: &web.Config{
			HTTPS:      false,
			CertFile:   "",
			KeyFile:    "",
			Port:       "8080",
			ServerName: "typing",
			Static: map[string]string{
				"/admin": "./static/admin/",
			},
		},

		DBDSN:    "./output/main.db",
		DBPrefix: "typing_",
		DBDriver: "sqlite3",

		FrontAPIPrefix: "/api",
		AdminAPIPrefix: "/admin/api",
		ThemeURLPrefix: "/themes",
		ThemeDir:       "./static/front/themes/",
		TempDir:        "./output/temp/",

		UploadDir:       "./output/uploads/",
		UploadExts:      ".txt;.png;.jpg;.jpeg",
		UploadSize:      1024 * 1024 * 5,
		UploadURLPrefix: "/uploads",
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, os.ModePerm)
}

// 向数据库写入初始内容。
func fillDB(db *orm.DB) error {
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

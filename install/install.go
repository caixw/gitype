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
	if err := ioutil.WriteFile(logsConfigPath, logFile, os.ModePerm); err != nil {
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
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	// option
	if err := fillOptions(db); err != nil {
		return err
	}

	// tags
	if err = fillTags(db); err != nil {
		return err
	}

	// post
	post := &models.Post{
		Title:    "第一篇日志",
		Content:  "<p>这是你的第一篇日志</p>",
		State:    models.PostStatePublished,
		Created:  time.Now().Unix(),
		Modified: time.Now().Unix(),
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

func fillOptions(db *orm.DB) error {
	opt := &core.Options{
		SiteURL:     "localhost:8081/",
		SiteName:    "typing blog",
		SecondTitle: "副标题",
		ScreenName:  "typing",
		Password:    core.HashPassword("123"),
		Keywords:    "typing",
		Description: "typing-极简的博客系统",
		Suffix:      ".html",

		PageSize:     20,
		DateFormat:   "2006-01-02 15:04:05",
		SidebarSize:  10,
		CommentOrder: core.CommentOrderDesc,

		PostsChangefreq: "never",
		TagsChangefreq:  "daily",
		PostsPriority:   0.9,
		TagsPriority:    0.4,
		RssSize:         20,

		Theme: "default",
	}

	maps, err := opt.ToMaps()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sql := "INSERT INTO #options ({key},{group},{value}) VALUES(?,?,?)"
	for _, item := range maps {
		_, err := tx.Exec(true, sql, item["key"], item["group"], item["value"])
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

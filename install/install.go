// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 用于执行typing的安装。
// 输出默认的配置文件；
// 输出默认的日志配置文件；
// 填充默认的数据到数据库；
package install

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

// OutputConfigFile 用于输出配置文件到指定的位置。
// 目前包含了日志配置文件和程序本身的配置文件。
func OutputConfigFile(logsConfigPath, configPath string) error {
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

		ThemeDir: "./static/front/",
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, os.ModePerm)
}

// FillDB 向数据库写入初始内容。
func FillDB(db *orm.DB) error {
	if db == nil {
		return errors.New("db==nil")
	}

	// option
	if err := db.Create(&models.Option{}); err != nil {
		return err
	}

	if err := fillOptions(db); err != nil {
		return err
	}

	// meta
	if err := db.Create(&models.Meta{}); err != nil {
		return err
	}
	metas := []*models.Meta{
		// cats
		{
			Name:        "default",
			Title:       "默认分类",
			Type:        models.MetaTypeCat,
			Order:       10,
			Parent:      models.MetaNoParent,
			Description: "所有添加的文章，默认添加此分类下。",
		},

		// tags
		{Name: "tag1", Title: "标签一", Type: models.MetaTypeTag, Description: "<h5>tag1</h5>"},
		{Name: "tag2", Title: "标签二", Type: models.MetaTypeTag, Description: "<h5>tag2</h5>"},
		{Name: "tag3", Title: "标签三", Type: models.MetaTypeTag, Description: "<h5>tag3</h5>"},
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
	if err := db.Create(&models.Post{}); err != nil {
		return err
	}

	if _, err := db.Insert(&models.Post{Title: "第一篇日志", Content: "<p>这是你的第一篇日志</p>"}); err != nil {
		return err
	}

	// models.Comment
	if err := db.Create(&models.Comment{}); err != nil {
		return err
	}

	if _, err := db.Insert(&models.Comment{PostID: 1, Content: "<p>沙发</p>", AuthorName: "游客"}); err != nil {
		return err
	}

	// relationship
	if err := db.Create(&models.Relationship{}); err != nil {
		return err
	}
	if _, err := db.Insert(&models.Relationship{MetaID: 1, PostID: 1}); err != nil {
		return err
	}

	return nil
}

func fillOptions(db *orm.DB) error {
	opt := &core.Options{
		PageSize:   20,
		SiteName:   "typing blog",
		ScreenName: "typing",
		Password:   core.HashPassword("123"),
		Theme:      "default",
		Keywords:   "typing",
	}

	maps, err := opt.ToMaps()
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

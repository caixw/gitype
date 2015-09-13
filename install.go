// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/caixw/typing/install"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

// 将默认的配置文件输出到`./config`目录下。
// 包含了`./config/logs.xml`和`./config/app.json`两个文件。
func outputConfig() error {
	ioutil.WriteFile(logConfigPath, install.LogFile, os.ModePerm)

	cfg := &config{
		Core: &web.Config{
			HTTPS:      false,
			CertFile:   "",
			KeyFile:    "",
			Port:       "8080",
			ServerName: "typing",
			Static: map[string]string{
				"/admin": "./static/admin/",
				"/":      "./static/front/",
				"":       "./static/front/",
			},
		},

		DBDSN:    "./output/main.db",
		DBPrefix: "typing_",
		DBDriver: "sqlite3",

		FrontAPIPrefix: "/api",
		AdminAPIPrefix: "/admin/api",
	}
	data, err := json.Marshal(cfg)
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
		{Name: "tag1", Title: "Tag1", Type: metaTypeTag, Description: "<h5>tag1</h5>"},
		{Name: "tag2", Title: "Tag2", Type: metaTypeTag, Description: "<h5>tag2</h5>"},
		{Name: "tag3", Title: "Tag3", Type: metaTypeTag, Description: "<h5>tag3</h5>"},

		{Name: "cat1", Title: "cat1", Type: metaTypeCat, Order: 10, Parent: -1, Description: "<h5>cat1</h5>"},
		{Name: "cat2", Title: "cat2", Type: metaTypeCat, Order: 10, Parent: -1, Description: "<h5>cat2</h5>"},
		{Name: "cat3", Title: "cat3", Type: metaTypeCat, Order: 10, Parent: -1, Description: "<h5>cat3</h5>"},
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

	// relationship
	if err := db.Create(&relationship{}); err != nil {
		return err
	}

	return nil
}

func fillOptions(db *orm.DB) error {
	opt := &options{
		Pretty:     true,
		PageSize:   20,
		SiteName:   "typing blog",
		ScreenName: "typing",
		Password:   hashPassword(defaultPassword),
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

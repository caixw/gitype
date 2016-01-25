// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/caixw/typing/app/static"
	"github.com/caixw/typing/models"
	"github.com/issue9/conv"
	"github.com/issue9/orm"
	"github.com/issue9/rands"
	"github.com/issue9/utils"
	"github.com/issue9/web"
)

// 执行安装程序。
func Install(appdir, action string) error {
	if !strings.HasSuffix(appdir, "/") && !strings.HasSuffix(appdir, string(os.PathSeparator)) {
		appdir += string(os.PathSeparator)
	}

	switch action {
	case "config":
		return installConfig(appdir)
	case "db":
		return installDB(appdir)
	default:
		return errors.New("app.Install:无效的action值")
	}
}

// 向数据库写入初始内容。
func installDB(appdir string) error {
	cfg, err := loadConfig(appdir + configFile)
	if err != nil {
		return err
	}

	db, err := initDB(cfg)
	if err != nil {
		return err
	}

	if err := models.Install(db); err != nil {
		return err
	}

	// option
	return fillOptions(db, cfg)
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
func fillOptions(db *orm.DB, cfg *Config) error {
	now := time.Now().Unix()
	opt := &Options{
		SiteName:    "typing blog",
		SecondTitle: "副标题",
		SiteURL:     "http://localhost:8080/",
		Keywords:    "typing",
		Description: "typing-极简的博客系统",
		Suffix:      ".html",
		Beian:       "备案中...",

		Uptime:      now,
		LastUpdated: now,
		Last:        "[]",

		PageSize:        20,
		SidebarSize:     10,
		LongDateFormat:  "2006-01-02 15:04:05",
		ShortDateFormat: "2006-01-02",
		CommentOrder:    CommentOrderDesc,

		PostsChangefreq: "never",
		TagsChangefreq:  "daily",
		PostsPriority:   0.9,
		TagsPriority:    0.4,
		RssSize:         20,

		Theme: "default",

		Menus: "[]",

		ScreenName: "typing",
		Email:      "",
		Password:   Password(defaultPassword),
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

// 用于输出配置文件到指定的位置。
// 目前包含了日志配置文件和程序本身的配置文件。
func installConfig(appdir string) error {
	if !utils.FileExists(appdir + configDir) {
		if err := os.MkdirAll(appdir+configDir, os.ModePerm); err != nil {
			return err
		}
		if !utils.FileExists(appdir + configDir) {
			return fmt.Errorf("路径[%v]不存在，且无法创建", appdir+configDir)
		}
	}

	if err := ioutil.WriteFile(appdir+logConfigFile, static.LogConfig, os.ModePerm); err != nil {
		return err
	}

	cfg := &Config{
		Core: &web.Config{
			HTTPS:    false,
			CertFile: "",
			KeyFile:  "",
			Port:     "8080",
			Headers: map[string]string{
				"Server": "typing",
			},
		},
		Debug: true,

		AdminURLPrefix: "/admin",
		AdminDir:       "./static/admin/",
		Salt:           rands.String(6, 7, rands.Lower, rands.Upper, rands.Digit, rands.Punct),

		DBDSN:    "./output/main.db",
		DBPrefix: "typing_",
		DBDriver: "sqlite3",

		FrontAPIPrefix: "/api",
		AdminAPIPrefix: "/admin/api",

		ThemeURLPrefix: "/themes",
		ThemeDir:       "./static/front/themes/",

		RootDir: "./static/front/root/",

		UploadDir:       "./output/uploads/",
		UploadDirFormat: "2006/01/",
		UploadExts:      ".txt;.png;.jpg;.jpeg",
		UploadSize:      1024 * 1024 * 5,
		UploadURLPrefix: "/uploads",
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(appdir+configFile, data, os.ModePerm)
}

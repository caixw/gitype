// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/forward"
	"github.com/issue9/web"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// 从app.json加载的配置内容。若要修改app.json的内容，需要重新运行程序。
// 只有重启才能配置的参数，写入到config中，其它参数可以在options中指定。
type config struct {
	Core *web.Config `json:"core"`

	DBDSN    string `json:"dbDSN"`
	DBPrefix string `json:"dbPrefix"`
	DBDriver string `json:"dbDriver"`

	FrontAPIPrefix string `json:"frontApiPrefix"` // 前端api地址前缀
	AdminAPIPrefix string `json:"adminApiPrefix"` // 后台api地址前经
}

// 加载配置文件到全局变量cfg中。
func loadConfig() (*config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// 检测AdminAPIPrefix是否符合要求
	if len(cfg.AdminAPIPrefix) == 0 {
		return nil, errors.New("必须指定adminApiPrefix值")
	}
	if cfg.AdminAPIPrefix[len(cfg.AdminAPIPrefix)-1] == '/' {
		return nil, errors.New("adminApiPrefix不能以/符号结尾")
	}

	// 检测FrontAPIPrefix是否符合要求
	if len(cfg.FrontAPIPrefix) == 0 {
		return nil, errors.New("必须指定frontApiPrefix值")
	}
	if cfg.FrontAPIPrefix[len(cfg.FrontAPIPrefix)-1] == '/' {
		return nil, errors.New("frontApiPrefix不能以/符号结尾")
	}

	if len(cfg.DBDSN) == 0 {
		return nil, errors.New("app.json中未指定dbDSN")
	}

	if len(cfg.DBDriver) == 0 {
		return nil, errors.New("app.json中未指定dbDriver")
	}

	return cfg, nil
}

// 从cfg中初始化一个*orm.DB实例
func initDB(cfg *config) (*orm.DB, error) {
	var d forward.Dialect
	switch cfg.DBDriver {
	case "sqlite3":
		d = dialect.Sqlite3()
	case "mysql":
		d = dialect.Mysql()
	case "postgres":
		d = dialect.Postgres()
	default:
		return nil, errors.New("不能理解的dbDriver值：" + cfg.DBDriver)
	}

	return orm.NewDB(cfg.DBDriver, cfg.DBDSN, cfg.DBPrefix, d)
}

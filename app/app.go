// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// app 定义了程序的初始化的基本内容，包括：
//  - 加载日志系统;
//  - 加载配置文件;
//  - 根据配置文件，初始化相应的数据库实例;
//  - 默认的配置文件安装脚本;
//  - 默认的日志配置文件安装脚本;
package app

import (
	"errors"

	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/forward"
)

const (
	Version = "0.9.41.151212" // 程序版本号

	// 定义两个配置文件的位置。
	configPath    = "./config/app.json"
	logConfigPath = "./config/logs.xml"
)

// 初始化系统，获取系统配置变量和数据库实例。
func Init() (*Config, *orm.DB, *Options, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, nil, nil, err
	}

	db, err := initDB(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = logs.InitFromXMLFile(logConfigPath); err != nil {
		return nil, nil, nil, err
	}

	opt, err := loadOptions(db)
	if err != nil {
		return nil, nil, nil, err
	}

	return cfg, db, opt, nil
}

// 从一个Config实例中初始一个orm.DB实例。
func initDB(cfg *Config) (*orm.DB, error) {
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

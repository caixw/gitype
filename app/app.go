// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// app 定义了程序的基本内容，包括：
//  - 加载日志系统;
//  - 加载配置文件;
//  - 根据配置文件，初始化相应的数据库实例;
//  - 从数据库加载配置内容及初始数据；
//  - 默认的配置文件安装脚本;
//  - 默认的日志配置文件安装脚本;
//  - 默信的数据库安装脚本;
package app

import (
	"errors"

	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/forward"
	"github.com/issue9/web"
)

const (
	Version = "0.13.68.160111" // 程序版本号

	defaultPassword = "123" // 默认的后台登录密码

	configPath    = "./config/app.json"
	logConfigPath = "./config/logs.xml"
)

type App struct {
	config  *Config
	options *Options
	stat    *Stat
	db      *orm.DB
}

// 初始化系统，获取系统配置变量和数据库实例。
func Init() (*App, error) {
	// 初始化日志系统
	if err := logs.InitFromXMLFile(logConfigPath); err != nil {
		return nil, err
	}

	// 加载app.json配置文件
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// 根据配置文件初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	// 加载数据库中的配置项
	opt, err := loadOptions(db)
	if err != nil {
		return nil, err
	}

	// 初始化系统的状态数据。
	stat, err := loadStat(db)
	if err != nil {
		return nil, err
	}

	return &App{
		config:  cfg,
		db:      db,
		options: opt,
		stat:    stat,
	}, nil
}

func (a *App) Config() *Config {
	return a.config
}

func (a *App) Options() *Options {
	return a.options
}

func (a *App) Stat() *Stat {
	return a.stat
}

func (a *App) DB() *orm.DB {
	return a.db
}

func (a *App) Run() {
	web.Run(a.config.Core)
}

func (a *App) Close() {
	a.db.Close()
	logs.Flush()
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

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
	"fmt"
	"os"
	"strings"

	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/forward"
	"github.com/issue9/utils"
	"github.com/issue9/web"
)

const (
	Version = "0.15.71.160119" // 程序版本号

	defaultPassword = "123" // 默认的后台登录密码

	configFile    = "config/app.json"
	logConfigFile = "config/logs.xml"
	adminDir      = "static/admin"
	themeDir      = "static/front/themes"
	rootDir       = "static/front/root"
	uploadDir     = "static/uploads"
)

var (
	appdir string

	config  *Config
	db      *orm.DB
	options *Options
	stats   *Stats
)

// 初始化app包。
// 除Install函数，其它函数都依赖Init()做初始化。
func Init(dir string) (err error) {
	if !utils.FileExists(dir) {
		return fmt.Errorf("appdir[%v]不存在", appdir)
	}

	if !strings.HasSuffix(dir, "/") && !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir += string(os.PathSeparator)
	}
	appdir = dir

	// 初始化日志系统
	if err := logs.InitFromXMLFile(Appdir(logConfigFile)); err != nil {
		return err
	}

	// 加载app.json配置文件
	config, err = loadConfig(Appdir(configFile))
	if err != nil {
		return err
	}

	// 根据配置文件初始化数据库
	db, err = initDB(config)
	if err != nil {
		return err
	}

	// 加载数据库中的配置项
	options, err = loadOptions()
	if err != nil {
		return err
	}

	// 初始化系统的状态数据。
	stats, err = loadStats()
	if err != nil {
		return err
	}

	return nil
}

func Appdir(path string) string {
	return appdir + path
}

func GetDB() *orm.DB {
	return db
}

func Run() {
	web.Run(config.Core)
}

func Close() {
	db.Close()
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

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/issue9/logs"
	"github.com/issue9/mux"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

const (
	version = "0.1.1.150914" // 版本号

	// 两个配置文件路径
	configPath    = "./config/app.json"
	logConfigPath = "./config/logs.xml"

	defaultPassword = "123" // 后台默认的登录密码
)

var (
	db  *orm.DB // 数据库实例
	opt *options
)

func main() {
	action := flag.String("init", "", "指定需要初始化的内容，可取的值可以为：config和db。")
	flag.Parse()
	switch *action {
	case "config":
		if err := outputConfig(); err != nil {
			panic(err)
		}
		return
	case "db":
		cfg, err := loadConfig()
		if err != nil {
			panic(err)
		}

		db, err := initDB(cfg)
		defer db.Close()
		if err != nil {
			panic(err)
		}
		if err := fillDB(db); err != nil {
			panic(err)
		}
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		panic(err)
	}

	db, err = initDB(cfg)
	if err != nil {
		panic(err)
	}

	if err := logs.InitFromXMLFile(logConfigPath); err != nil {
		panic(err)
	}

	if opt, err = loadOptions(); err != nil {
		panic(err)
	}

	if err := initThemes(cfg.ThemeDir); err != nil {
		panic(err)
	}

	if err := initModule(cfg); err != nil {
		panic(err)
	}

	cfg.Core.ErrHandler = mux.PrintDebug
	web.Run(cfg.Core)
	db.Close()
}

// 初始化模块，及与模块相对应的路由。
func initModule(cfg *config) error {
	m, err := web.NewModule("admin")
	if err != nil {
		return err
	}
	initAdminAPIRoutes(m.Prefix(cfg.AdminAPIPrefix))

	m, err = web.NewModule("front")
	if err != nil {
		return err
	}
	initFrontAPIRoutes(m.Prefix(cfg.FrontAPIPrefix))

	initFrontPageRoutes(m)

	return nil
}

func initFrontPageRoutes(m *web.Module) {
	m.GetFunc("/", pageIndex)
}

func initFrontAPIRoutes(front *mux.Prefix) {
	front.GetFunc("/tags", frontGetTags).
		GetFunc("/cats", frontGetCats)

	// post
	front.PostFunc("/posts/{id:\\d+}/comments", frontPostPostComment).
		GetFunc("/posts/{id:\\d+}", frontGetPost).
		GetFunc("/posts/{id:\\d+}/comments", frontGetPostComments).
		GetFunc("/posts", frontGetPosts)
}

func initAdminAPIRoutes(admin *mux.Prefix) {
	admin.PostFunc("/login", adminPostLogin).
		Delete("/login", loginHandlerFunc(adminDeleteLogin)).
		Put("/password", loginHandlerFunc(adminChangePassword))

	// options
	admin.Get("/options/{key}", loginHandlerFunc(adminGetOption)).
		Patch("/options/{key}", loginHandlerFunc(adminPatchOption))

	// cats
	admin.Put("/cats/{id:\\d+}", loginHandlerFunc(adminPutCat)).
		Delete("/cats/{id:\\d+}", loginHandlerFunc(adminDeleteCat)).
		Post("/cats", loginHandlerFunc(adminPostCat))

	// tags
	admin.Put("/tags/{id:\\d+}", loginHandlerFunc(adminPutTag)).
		Delete("/tags/{id:\\d+}", loginHandlerFunc(adminDeleteTag)).
		Post("/tags", loginHandlerFunc(adminPostTag))

	// comments
	admin.Get("/comments", loginHandlerFunc(getComments)).
		Post("/comments", loginHandlerFunc(adminPostComment)).
		Put("/comments/{id:\\d+}", loginHandlerFunc(putComment)).
		Post("/comments/{id:\\d+}/waiting", loginHandlerFunc(setCommentWaiting)).
		Post("/comments/{id:\\d+}/spam", loginHandlerFunc(setCommentSpam)).
		Post("/comments/{id:\\d+}/approved", loginHandlerFunc(setCommentApproved))

	admin.Get("/posts", loginHandlerFunc(adminGetPosts)).
		Post("/posts", loginHandlerFunc(adminPostPost)).
		Get("/posts/{id:\\d+}", loginHandlerFunc(adminGetPost)).
		Delete("/posts/{id:\\d+}", loginHandlerFunc(adminDeletePost)).
		Put("/posts/{id:\\d+}", loginHandlerFunc(adminPutPost))
}

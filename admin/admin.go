// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"
	"strings"
	"time"

	"github.com/caixw/typing/app"
	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/upload"
	"github.com/issue9/web"
)

var (
	cfg  *app.Config
	db   *orm.DB
	opt  *app.Options
	stat *app.Stat
	u    *upload.Upload
)

// 初始化当前模块
func Init(config *app.Config, database *orm.DB, options *app.Options, s *app.Stat) error {
	cfg = config
	opt = options
	db = database
	stat = s

	// 上传相关配置
	var err error
	u, err = upload.New(cfg.UploadDir, cfg.UploadDirFormat, cfg.UploadSize, strings.Split(cfg.UploadExts, ";")...)
	if err != nil {
		return err
	}

	return initRoute()
}

// 更新数据库中的lastUpdated变量
func lastUpdated() {
	if err := opt.Set(db, "lastUpdated", time.Now().Unix(), true); err != nil {
		logs.Error("admin.lastUpdated:", err)
	}
}

func initRoute() error {
	m, err := web.NewModule("admin")
	if err != nil {
		return err
	}

	m.Get(cfg.AdminURLPrefix+"/", http.StripPrefix(cfg.AdminURLPrefix, http.FileServer(http.Dir(cfg.AdminDir))))

	p := m.Prefix(cfg.AdminAPIPrefix)

	p.Get("/state", loginHandlerFunc(adminGetState)).
		Get("/modules", loginHandlerFunc(adminGetModules)).
		Put("/modules/{name}/start", loginHandlerFunc(adminPutModuleStart)).
		Put("/modules/{name}/stop", loginHandlerFunc(adminPutModuleStop))

	// users
	p.PostFunc("/login", adminPostLogin).
		Delete("/login", loginHandlerFunc(adminDeleteLogin)).
		Put("/password", loginHandlerFunc(adminChangePassword))

	// feed
	p.Put("/feed/sitemap", loginHandlerFunc(adminPutSitemap)).
		Put("/feed/rss", loginHandlerFunc(adminPutRss)).
		Put("/feed/atom", loginHandlerFunc(adminPutAtom))

	// meida
	p.Post("/media", loginHandlerFunc(adminPostMedia)).
		Get("/media", loginHandlerFunc(adminGetMedia))

	// themes
	p.Get("/themes", loginHandlerFunc(adminGetThemes)).
		Get("/themes/current", loginHandlerFunc(adminGetCurrentTheme)).
		Put("/themes/current", loginHandlerFunc(adminPutCurrentTheme))

	// options
	p.Get("/options/{key}", loginHandlerFunc(adminGetOption)).
		Patch("/options/{key}", loginHandlerFunc(adminPatchOption))

	// tags
	p.Put("/tags/{id:\\d+}", loginHandlerFunc(adminPutTag)).
		Delete("/tags/{id:\\d+}", loginHandlerFunc(adminDeleteTag)).
		Get("/tags/{id:\\d+}", loginHandlerFunc(adminGetTag)).
		Post("/tags", loginHandlerFunc(adminPostTag)).
		Get("/tags", loginHandlerFunc(adminGetTags))

	// comments
	p.Get("/comments", loginHandlerFunc(adminGetComments)).
		Get("/comments/count", loginHandlerFunc(adminGetCommentsCount)).
		Post("/comments", loginHandlerFunc(adminPostComment)).
		Put("/comments/{id:\\d+}", loginHandlerFunc(adminPutComment)).
		Delete("/comments/{id:\\d+}", loginHandlerFunc(adminDeleteComment)).
		Post("/comments/{id:\\d+}/waiting", loginHandlerFunc(adminSetCommentWaiting)).
		Post("/comments/{id:\\d+}/spam", loginHandlerFunc(adminSetCommentSpam)).
		Post("/comments/{id:\\d+}/approved", loginHandlerFunc(adminSetCommentApproved))

	// posts
	p.Get("/posts", loginHandlerFunc(adminGetPosts)).
		Get("/posts/count", loginHandlerFunc(adminGetPostsCount)).
		Post("/posts", loginHandlerFunc(adminPostPost)).
		Get("/posts/{id:\\d+}", loginHandlerFunc(adminGetPost)).
		Delete("/posts/{id:\\d+}", loginHandlerFunc(adminDeletePost)).
		Put("/posts/{id:\\d+}", loginHandlerFunc(adminPutPost)).
		Post("/posts/{id:\\d+}/draft", loginHandlerFunc(adminSetPostDraft)).
		Post("/posts/{id:\\d+}/published", loginHandlerFunc(adminSetPostPublished))

	return nil
}

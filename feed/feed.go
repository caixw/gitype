// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/feed/static"
	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

var (
	db  *orm.DB
	opt *app.Options

	sitemap      = new(bytes.Buffer)
	sitemapMutex sync.Mutex

	atom      = new(bytes.Buffer)
	atomMutex sync.Mutex

	rss      = new(bytes.Buffer)
	rssMutex sync.Mutex
)

// 初始化sitemap包，path为sitemap.xml文件的保存路径
func Init(database *orm.DB, options *app.Options) error {
	db = database
	opt = options

	if err := BuildRss(); err != nil {
		return err
	}
	if err := BuildAtom(); err != nil {
		return err
	}
	if err := BuildSitemap(); err != nil {
		return err
	}

	return initRoute()
}

// 初始化路由项
func initRoute() error {
	m, err := web.NewModule("feed")
	if err != nil {
		return err
	}

	m.GetFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		sitemapMutex.Lock()
		defer sitemapMutex.Unlock()

		if _, err := w.Write(sitemap.Bytes()); err != nil {
			logs.Error("feed.initRoute.route-/sitemap.xml:", err)
			w.WriteHeader(404) // 若是出错，给客户端的信息提示为404
		}
	})

	// NOTE:若修改此路由，请同时修改sitemap.xml中的相对应的.xsl路径
	m.GetFunc("/sitemap.xsl", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(static.Sitemap); err != nil {
			logs.Error("feed.initRoute.route-/sitemap.xsl:", err)
			w.WriteHeader(404)
		}
	})

	m.GetFunc("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		rssMutex.Lock()
		defer rssMutex.Unlock()

		if _, err := w.Write(rss.Bytes()); err != nil {
			logs.Error("feed.initRoute.route-/rss.xml:", err)
			w.WriteHeader(404)
		}
	})

	m.GetFunc("/atom.xml", func(w http.ResponseWriter, r *http.Request) {
		atomMutex.Lock()
		defer atomMutex.Unlock()

		if _, err := w.Write(atom.Bytes()); err != nil {
			logs.Error("feed.initRoute.route-/atom.xml:", err)
			w.WriteHeader(404)
		}
	})

	return nil
}

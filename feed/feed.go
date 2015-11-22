// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"net/http"
	"os"

	"github.com/caixw/typing/core"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

// 定义了各文件名。
const (
	sitemap    = "sitemap.xml"
	sitemapXsl = "sitemap.xsl"
	rss        = "rss.xml"
	atom       = "atom.xml"
)

var (
	db  *orm.DB
	opt *core.Options

	sitemapPath    string
	sitemapXslPath string
	rssPath        string
	atomPath       string
)

// 初始化sitemap包，path为sitemap.xml文件的保存路径
func Init() error {
	sitemapPath = core.Cfg.TempDir + sitemap
	sitemapXslPath = core.Cfg.TempDir + sitemapXsl
	rssPath = core.Cfg.TempDir + rss
	atomPath = core.Cfg.TempDir + atom
	db = core.DB
	opt = core.Opt

	file, err := os.Create(sitemapXslPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(files); err != nil {
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

	m.GetFunc("/"+sitemap, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapPath)
	})
	m.GetFunc("/"+sitemapXsl, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapXslPath)
	})
	m.GetFunc("/"+rss, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rssPath)
	})
	m.GetFunc("/"+atom, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, atomPath)
	})

	return nil
}

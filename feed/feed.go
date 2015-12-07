// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"net/http"
	"os"

	"github.com/caixw/typing/boot"
	"github.com/caixw/typing/feed/static"
	"github.com/caixw/typing/options"
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
	opt *options.Options

	sitemapPath    string
	sitemapXslPath string
	rssPath        string
	atomPath       string
)

// 初始化sitemap包，path为sitemap.xml文件的保存路径
func Init(cfg *boot.Config, database *orm.DB, options *options.Options) error {
	sitemapPath = cfg.TempDir + sitemap
	sitemapXslPath = cfg.TempDir + sitemapXsl
	rssPath = cfg.TempDir + rss
	atomPath = cfg.TempDir + atom
	db = database
	opt = options

	// 输出sitemap.xsl到临时目录
	file, err := os.Create(sitemapXslPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(static.Sitemap); err != nil {
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

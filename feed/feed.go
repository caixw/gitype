// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"net/http"

	"github.com/caixw/typing/core"
	"github.com/issue9/orm"
	"github.com/issue9/web"
)

const (
	sitemap     = "sitemap.xml"
	sitemapXslt = "sitemap.xslt"
	rss         = "rss.xml"
	atom        = "atom.xml"
)

var (
	db  *orm.DB
	opt *core.Options

	sitemapPath     string
	sitemapXsltPath string
	rssPath         string
	atomPath        string
)

// 初始化sitemap包，path为sitemap.xml文件的保存路径
func Init(tempDir string, database *orm.DB, opt *core.Options) {
	sitemapPath = tempDir + sitemap
	sitemapXsltPath = tempDir + sitemapXslt
	rssPath = tempDir + rss
	atomPath = tempDir + atom
	db = database
}

func InitRoute(w *web.Module) {
	w.GetFunc("/"+sitemap, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapPath)
	})

	w.GetFunc("/"+sitemapXslt, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapXsltPath)
	})

	w.GetFunc("/"+rss, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rssPath)
	})

	w.GetFunc("/"+atom, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, atomPath)
	})
}

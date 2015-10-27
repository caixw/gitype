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
func Init(tempDir string, database *orm.DB, options *core.Options) error {
	sitemapPath = tempDir + sitemap
	sitemapXslPath = tempDir + sitemapXsl
	rssPath = tempDir + rss
	atomPath = tempDir + atom
	db = database
	opt = options

	file, err := os.Create(sitemapXslPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(files)
	return err
}

// 初始化路由项
func InitRoute(w *web.Module) {
	w.GetFunc("/"+sitemap, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapPath)
	})

	w.GetFunc("/"+sitemapXsl, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, sitemapXslPath)
	})

	w.GetFunc("/"+rss, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, rssPath)
	})

	w.GetFunc("/"+atom, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, atomPath)
	})
}

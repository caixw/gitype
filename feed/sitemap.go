// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm"
	"github.com/issue9/orm/fetch"
)

const (
	sitemapHeader = `<?xml version="1.0" encoding="utf-8"?>
<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	sitemapFooter = `</urlset>`
)

// Build 构建一个sitemap.xml文件到sitemapPath文件中，若该文件已经存在，则覆盖。
func BuildSitemap() error {
	buf := bytes.NewBufferString(sitemapHeader)
	buf.Grow(10000)

	if err := addPostsToSitemap(buf, db, opt); err != nil {
		return err
	}

	if err := addTagsToSitemap(buf, db, opt); err != nil {
		return err
	}

	buf.WriteString(sitemapFooter)

	file, err := os.Create(sitemapPath)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(file)
	file.Close()
	return err
}

func addPostsToSitemap(buf *bytes.Buffer, db *orm.DB, opt *core.Options) error {
	sql := "SELECT {id}, {name}, {title}, {summary}, {content}, {created}, {modified} FROM #posts WHERE {state}=?"
	rows, err := db.Query(true, sql, models.PostStatePublished)
	if err != nil {
		return err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return err
	}

	for _, v := range maps {
		loc := opt.SiteURL + "/posts/"
		if len(v["name"]) > 0 {
			loc += v["name"]
		} else {
			loc += v["id"]
		}

		modified, err := strconv.ParseInt(v["modified"], 10, 64)
		if err != nil {
			return err
		}
		addItemToSitemap(buf, loc, opt.PostsChangefreq, modified, opt.PostsPriority)
	}
	return nil
}

func addTagsToSitemap(buf *bytes.Buffer, db *orm.DB, opt *core.Options) error {
	sql := "SELECT {id}, {name}, {title}, {description} FROM #tags"
	rows, err := db.Query(true, sql)
	if err != nil {
		return err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return err
	}

	for _, v := range maps {
		loc := opt.SiteURL + "/tags/" + v["name"]

		addItemToSitemap(buf, loc, opt.TagsChangefreq, time.Now().Unix(), opt.TagsPriority)
	}
	return nil
}

func addItemToSitemap(buf *bytes.Buffer, loc, changefreq string, lastmod int64, priority float64) {
	buf.WriteString("<url>\n")

	buf.WriteString("<loc>")
	buf.WriteString(loc)
	buf.WriteString("</loc>\n")

	t := time.Unix(lastmod, 0)
	buf.WriteString("<lastmod>")
	buf.WriteString(t.Format("2006-01-02T15:04:05+08:00"))
	buf.WriteString("</lastmod>\n")

	buf.WriteString("<changefreq>")
	buf.WriteString(changefreq)
	buf.WriteString("</changefreq>\n")

	buf.WriteString("<priority>")
	buf.WriteString(strconv.FormatFloat(priority, 'f', 1, 32))
	buf.WriteString("</priority>\n")

	buf.WriteString("</url>\n")
}

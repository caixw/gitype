// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"strconv"
	"time"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/themes"
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
	if _, err := sitemapW.WriteString(sitemapHeader); err != nil {
		return err
	}

	if err := addPostsToSitemap(sitemapW, db, opt); err != nil {
		return err
	}

	if err := addTagsToSitemap(sitemapW, db, opt); err != nil {
		return err
	}

	if _, err := sitemapW.WriteString(sitemapFooter); err != nil {
		return err
	}

	sitemapR, sitemapW = sitemapW, sitemapR
	sitemapW.Reset()
	return nil
}

func addPostsToSitemap(buf *bytes.Buffer, db *orm.DB, opt *app.Options) error {
	sql := `SELECT {id} AS ID, {name} AS Name, {title} AS Title, {summary} AS Summary,
		{content} AS Content, {created} AS Created, {modified} AS Modified
		FROM #posts WHERE {state}=?`
	rows, err := db.Query(true, sql, models.PostStatePublished)
	if err != nil {
		return err
	}
	defer rows.Close()

	posts := make([]*themes.Post, 0, 100)
	if _, err := fetch.Obj(&posts, rows); err != nil {
		return err
	}

	for _, p := range posts {
		addItemToSitemap(buf, p.Permalink(), opt.PostsChangefreq, p.Modified, opt.PostsPriority)
	}
	return nil
}

func addTagsToSitemap(buf *bytes.Buffer, db *orm.DB, opt *app.Options) error {
	sql := "SELECT {id} AS ID, {name} AS Name, {title} AS Title, {description} AS Description FROM #tags"
	rows, err := db.Query(true, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	tags := make([]*themes.Tag, 0, 100)
	if _, err := fetch.Obj(&tags, rows); err != nil {
		return err
	}

	now := time.Now().Unix()
	for _, tag := range tags {
		addItemToSitemap(buf, tag.Permalink(), opt.TagsChangefreq, now, opt.TagsPriority)
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

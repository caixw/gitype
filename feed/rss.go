// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"time"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/themes"
	"github.com/issue9/orm"
	"github.com/issue9/orm/fetch"
)

const (
	rssHeader = `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
<channel>`

	rssFooter = `</channel>
</rss>`
)

// Build 构建一个rss.xml文件到rssPath文件中，若该文件已经存在，则覆盖。
func BuildRss() error {
	if _, err := rssW.WriteString(rssHeader); err != nil {
		return err
	}

	if err := addPostsToRss(rssW, db, opt); err != nil {
		return err
	}

	if _, err := rssW.WriteString(rssFooter); err != nil {
		return err
	}

	rssR, rssW = rssW, rssR
	rssW.Reset()
	return nil
}

func addPostsToRss(buf *bytes.Buffer, db *orm.DB, opt *app.Options) error {
	sql := `SELECT {id} AS ID, {name} AS Name, {title} AS Title, {summary} AS Summary,
		{content} AS Content, {created} AS Created, {modified} AS Modified
		FROM #posts WHERE {state}=? LIMIT ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, opt.RssSize)
	if err != nil {
		return err
	}
	defer rows.Close()

	posts := make([]*themes.Post, 0, 100)
	if _, err := fetch.Obj(&posts, rows); err != nil {
		return err
	}

	for _, p := range posts {
		addItemToRss(buf, p.Permalink(), p.Title, p.Entry(), p.Modified)
	}
	return nil
}

func addItemToRss(buf *bytes.Buffer, link, title, description string, pubDate int64) {
	buf.WriteString("<item>\n")

	buf.WriteString("<link>")
	buf.WriteString(link)
	buf.WriteString("</link>\n")

	buf.WriteString("<title>")
	buf.WriteString(title)
	buf.WriteString("</title>\n")

	t := time.Unix(pubDate, 0)
	buf.WriteString("<pubDate>")
	buf.WriteString(t.Format(time.RFC1123))
	buf.WriteString("</pubDate>\n")

	buf.WriteString("<description>")
	buf.WriteString(description)
	buf.WriteString("</description>\n")

	buf.WriteString("</item>\n")
}

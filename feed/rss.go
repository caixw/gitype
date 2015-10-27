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
	rssHeader = `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
<channel>`

	rssFooter = `</channel>
</rss>`
)

// Build 构建一个sitemap.xml文件到sitemapPath文件中，若该文件已经存在，则覆盖。
func BuildRss() error {
	buf := bytes.NewBufferString(rssHeader)
	buf.Grow(10000)

	if err := addPostsToRss(buf, db, opt); err != nil {
		return err
	}

	buf.WriteString(rssFooter)

	file, err := os.Create(sitemapPath)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(file)
	file.Close()
	return err
}

func addPostsToRss(buf *bytes.Buffer, db *orm.DB, opt *core.Options) error {
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
		link := opt.SiteURL + "/posts/"
		if len(v["name"]) > 0 {
			link += v["name"]
		} else {
			link += v["id"]
		}

		modified, err := strconv.ParseInt(v["modified"], 10, 64)
		if err != nil {
			return err
		}
		addItemToRss(buf, link, v["title"], v["summary"], modified)
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

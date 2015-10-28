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
	atomHeader = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom>"`

	atomFooter = `</feed>`
)

// Build 构建一个atom.xml文件到atomPath文件中，若该文件已经存在，则覆盖。
func BuildAtom() error {
	buf := bytes.NewBufferString(atomHeader)
	buf.Grow(10000)

	buf.WriteString("<id>")
	buf.WriteString(opt.SiteURL)
	buf.WriteString("</id>\n")

	buf.WriteString("<link>")
	buf.WriteString(opt.SiteURL)
	buf.WriteString("</link>\n")

	buf.WriteString("<title>")
	buf.WriteString(opt.SiteName)
	buf.WriteString("</title>\n")

	buf.WriteString("<subtitle>")
	buf.WriteString(opt.SecondTitle)
	buf.WriteString("</subtitle>\n")

	buf.WriteString("<update>")
	buf.WriteString(time.Now().Format("2006-01-02T15:04:05Z08:00"))
	buf.WriteString("</update>\n")

	if err := addPostsToRss(buf, db, opt); err != nil {
		return err
	}

	buf.WriteString(atomFooter)

	file, err := os.Create(atomPath)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(file)
	file.Close()
	return err
}

func addPostsToAtom(buf *bytes.Buffer, db *orm.DB, opt *core.Options) error {
	sql := `SELECT {id}, {name}, {title}, {summary}, {content}, {created}, {modified}
	FROM #posts
	WHERE {state}=?
	LIMIT ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, opt.RssSize)
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

func addItemToAtom(buf *bytes.Buffer, link, title, summary string, update int64) {
	buf.WriteString("<entry>\n")

	buf.WriteString("<id>")
	buf.WriteString(link)
	buf.WriteString("</id>\n")

	buf.WriteString("<link>")
	buf.WriteString(link)
	buf.WriteString("</link>\n")

	buf.WriteString("<title>")
	buf.WriteString(title)
	buf.WriteString("</title>\n")

	t := time.Unix(update, 0)
	buf.WriteString("<update>")
	buf.WriteString(t.Format("2006-01-02T15:04:05Z08:00"))
	buf.WriteString("</update>\n")

	buf.WriteString("<summary>")
	buf.WriteString(summary)
	buf.WriteString("</summary>\n")

	buf.WriteString("</entry>\n")
}

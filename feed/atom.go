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
	atomHeader = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom>"`

	atomFooter = `</feed>`
)

// Build 构建一个atom.xml文件到atomPath文件中，若该文件已经存在，则覆盖。
func BuildAtom() error {
	if _, err := atomW.WriteString(atomHeader); err != nil {
		return err
	}

	atomW.WriteString("<id>")
	atomW.WriteString(opt.SiteURL)
	atomW.WriteString("</id>\n")

	atomW.WriteString("<link>")
	atomW.WriteString(opt.SiteURL)
	atomW.WriteString("</link>\n")

	atomW.WriteString("<title>")
	atomW.WriteString(opt.SiteName)
	atomW.WriteString("</title>\n")

	atomW.WriteString("<subtitle>")
	atomW.WriteString(opt.SecondTitle)
	atomW.WriteString("</subtitle>\n")

	atomW.WriteString("<update>")
	atomW.WriteString(time.Now().Format("2006-01-02T15:04:05Z08:00"))
	atomW.WriteString("</update>\n")

	if err := addPostsToRss(atomW, db, opt); err != nil {
		return err
	}

	if _, err := atomW.WriteString(atomFooter); err != nil {
		return err
	}

	atomR, atomW = atomW, atomR
	atomW.Reset()
	return nil
}

func addPostsToAtom(buf *bytes.Buffer, db *orm.DB, opt *app.Options) error {
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
		addItemToAtom(buf, p.Permalink(), p.Title, p.Entry(), p.Modified)
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

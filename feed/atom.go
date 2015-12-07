// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"time"

	"github.com/caixw/typing/models"
	"github.com/caixw/typing/options"
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
	atom.Reset()

	if _, err := atom.WriteString(atomHeader); err != nil {
		return err
	}

	atom.WriteString("<id>")
	atom.WriteString(opt.SiteURL)
	atom.WriteString("</id>\n")

	atom.WriteString("<link>")
	atom.WriteString(opt.SiteURL)
	atom.WriteString("</link>\n")

	atom.WriteString("<title>")
	atom.WriteString(opt.SiteName)
	atom.WriteString("</title>\n")

	atom.WriteString("<subtitle>")
	atom.WriteString(opt.SecondTitle)
	atom.WriteString("</subtitle>\n")

	atom.WriteString("<update>")
	atom.WriteString(time.Now().Format("2006-01-02T15:04:05Z08:00"))
	atom.WriteString("</update>\n")

	if err := addPostsToRss(atom, db, opt); err != nil {
		return err
	}

	_, err := atom.WriteString(atomFooter)

	return err
}

func addPostsToAtom(buf *bytes.Buffer, db *orm.DB, opt *options.Options) error {
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

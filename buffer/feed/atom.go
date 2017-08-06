// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"time"

	"github.com/caixw/typing/data"
)

const (
	atomHeader = `<feed xmlns="http://www.w3.org/2005/Atom"
      xmlns:opensearch="http://a9.com/-/spec/opensearch/1.1/">`

	atomFooter = `</feed>`
)

// BuildAtom 用于生成一个符合 atom 规范的 XML 文本 buffer。
func BuildAtom(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := &errWriter{
		buf: buf,
	}

	w.writeString(xmlHeader)

	w.writeString(atomHeader)

	w.writeString("<id>")
	w.writeString(d.Config.URL)
	w.writeString("</id>\n")

	w.writeString(`<link href="`)
	w.writeString(d.Config.URL)
	w.writeString("\" />\n")

	if d.Config.Opensearch != nil {
		o := d.Config.Opensearch
		w.writeString(`<link rel="search" type="application/opensearchdescription+xml" href="`)
		w.writeString(d.Config.URL + o.URL)
		w.writeString(`" title="`)
		w.writeString(o.Title)
		w.writeString("\" />\n")
	}

	w.writeString("<title>")
	w.writeString(d.Config.Title)
	w.writeString("</title>\n")

	w.writeString("<subtitle>")
	w.writeString(d.Config.Subtitle)
	w.writeString("</subtitle>\n")

	w.writeString("<update>")
	w.writeString(time.Now().Format("2006-01-02T15:04:05Z07:00"))
	w.writeString("</update>\n")

	addPostsToAtom(w, d)

	w.writeString(atomFooter)

	if w.err != nil {
		return nil, w.err
	}
	return buf, nil
}

func addPostsToAtom(w *errWriter, d *data.Data) {
	for _, p := range d.Posts {
		w.writeString("<entry>\n")

		w.writeString("<id>")
		w.writeString(p.Permalink)
		w.writeString("</id>\n")

		w.writeString(`<link href="`)
		w.writeString(d.Config.URL + p.Permalink)
		w.writeString("\" />\n")

		w.writeString("<title>")
		w.writeString(p.Title)
		w.writeString("</title>\n")

		t := time.Unix(p.Modified, 0)
		w.writeString("<update>")
		w.writeString(t.Format("2006-01-02T15:04:05Z07:00"))
		w.writeString("</update>\n")

		w.writeString("<summary>")
		w.writeString(p.Summary)
		w.writeString("</summary>\n")

		w.writeString("</entry>\n")
	}
}

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
	rssHeader = `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>`

	rssFooter = `</channel>
</rss>`
)

// BuildRSS 生成一个符合 rss 规范的 XML 文本 buffer。
func BuildRSS(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := &errWriter{
		buf: buf,
	}

	w.writeString(xmlHeader)

	w.writeString(rssHeader)

	w.writeString("\n<title>")
	w.writeString(d.Config.Title)
	w.writeString("</title>\n")

	w.writeString("<description>")
	w.writeString(d.Config.Subtitle)
	w.writeString("</description>\n")

	w.writeString("<link>")
	w.writeString(d.Config.URL)
	w.writeString("</link>\n")

	if d.Config.Opensearch != nil {
		o := d.Config.Opensearch
		w.writeString(`<atom:link rel="search" type="application/opensearchdescription+xml" href="`)
		w.writeString(d.Config.URL + o.URL)
		w.writeString(`" title="`)
		w.writeString(o.Title)
		w.writeString("\" />\n")
	}

	addPostsToRSS(w, d)

	w.writeString(rssFooter)

	if w.err != nil {
		return nil, w.err
	}
	return buf, nil
}

func addPostsToRSS(w *errWriter, d *data.Data) {
	for _, p := range d.Posts {
		w.writeString("<item>\n")

		w.writeString("<link>")
		w.writeString(d.Config.URL + p.Permalink)
		w.writeString("</link>\n")

		w.writeString("<title>")
		w.writeString(p.Title)
		w.writeString("</title>\n")

		t := time.Unix(p.Created, 0)
		w.writeString("<pubDate>")
		w.writeString(t.Format(time.RFC1123))
		w.writeString("</pubDate>\n")

		w.writeString("<description>")
		w.writeString(p.Summary)
		w.writeString("</description>\n")

		w.writeString("</item>\n")
	}
}

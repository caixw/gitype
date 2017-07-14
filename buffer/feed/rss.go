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
	rssHeader = `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>`

	rssFooter = `</channel>
</rss>`
)

// BuildRSS 生成一个符合 rss 规范的 XML 文本 buffer。
func BuildRSS(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(rssHeader); err != nil {
		return nil, err
	}

	buf.WriteString("\n<title>")
	buf.WriteString(d.Config.Title)
	buf.WriteString("</title>\n")

	buf.WriteString("<description>")
	buf.WriteString(d.Config.Subtitle)
	buf.WriteString("</description>\n")

	buf.WriteString("<link>")
	buf.WriteString(d.Config.URL)
	buf.WriteString("</link>\n")

	if d.Config.Opensearch != nil {
		o := d.Config.Opensearch
		buf.WriteString(`<atom:link rel="search" type="application/opensearchdescription+xml" href="`)
		buf.WriteString(d.Config.URL + o.URL)
		buf.WriteString(`" title="`)
		buf.WriteString(o.Title)
		buf.WriteString("\" />\n")
	}

	addPostsToRSS(buf, d)

	if _, err := buf.WriteString(rssFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

func addPostsToRSS(buf *bytes.Buffer, d *data.Data) {
	for _, p := range d.Posts {
		buf.WriteString("<item>\n")

		buf.WriteString("<link>")
		buf.WriteString(d.Config.URL + p.Permalink)
		buf.WriteString("</link>\n")

		buf.WriteString("<title>")
		buf.WriteString(p.Title)
		buf.WriteString("</title>\n")

		t := time.Unix(p.Created, 0)
		buf.WriteString("<pubDate>")
		buf.WriteString(t.Format(time.RFC1123))
		buf.WriteString("</pubDate>\n")

		buf.WriteString("<description>")
		buf.WriteString(p.Summary)
		buf.WriteString("</description>\n")

		buf.WriteString("</item>\n")
	}
}

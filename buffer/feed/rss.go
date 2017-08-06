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
	w := &writer{
		buf: buf,
	}

	w.writeString(xmlHeader)
	w.writeString(rssHeader)

	w.writeElement("title", d.Config.Title, nil)
	w.writeElement("description", d.Config.Subtitle, nil)
	w.writeElement("link", d.Config.URL, nil)

	if d.Config.Opensearch != nil {
		o := d.Config.Opensearch

		w.writeCloseElement("atom:link", map[string]string{
			"rel":   "search",
			"type":  "application/opensearchdescription+xml",
			"title": o.Title,
			"href":  d.Config.URL + o.URL,
		})
	}

	addPostsToRSS(w, d)

	w.writeString(rssFooter)

	if w.err != nil {
		return nil, w.err
	}
	return buf, nil
}

func addPostsToRSS(w *writer, d *data.Data) {
	for _, p := range d.Posts {
		w.writeString("<item>\n")

		w.writeElement("link", d.Config.URL+p.Permalink, nil)
		w.writeElement("title", p.Title, nil)
		t := time.Unix(p.Created, 0)
		w.writeElement("pubDate", t.Format(time.RFC1123), nil)
		w.writeElement("description", p.Summary, nil)

		w.writeString("</item>\n")
	}
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"time"

	"github.com/caixw/typing/data"
)

// BuildRSS 生成一个符合 rss 规范的 XML 文本 buffer。
func BuildRSS(d *data.Data) ([]byte, error) {
	w := newWrite()

	w.writeStartElement("rss", map[string]string{
		"version":    "2.0",
		"xmlns:atom": "http://www.w3.org/2005/Atom",
	}, true)
	w.writeStartElement("channel", nil, true)

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

	w.writeEndElement("channel")
	w.writeEndElement("rss")

	return w.bytes()
}

func addPostsToRSS(w *writer, d *data.Data) {
	for _, p := range d.Posts {
		w.writeStartElement("item", nil, true)

		w.writeElement("link", d.Config.URL+p.Permalink, nil)
		w.writeElement("title", p.Title, nil)
		t := time.Unix(p.Created, 0)
		w.writeElement("pubDate", t.Format(time.RFC1123), nil)
		w.writeElement("description", p.Summary, nil)

		w.writeEndElement("item")
	}
}

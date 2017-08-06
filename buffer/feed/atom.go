// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"time"

	"github.com/caixw/typing/data"
)

// BuildAtom 用于生成一个符合 atom 规范的 XML 文本 buffer。
func BuildAtom(d *data.Data) ([]byte, error) {
	w := newWrite()

	w.writeStartElement("feed", map[string]string{
		"xmlns":            "http://www.w3.org/2005/Atom",
		"xmlns:opensearch": "http://a9.com/-/spec/opensearch/1.1/",
	}, true)

	w.writeElement("id", d.Config.URL, nil)

	w.writeCloseElement("link", map[string]string{
		"href": d.Config.URL,
	})

	if d.Config.Opensearch != nil {
		o := d.Config.Opensearch

		w.writeCloseElement("link", map[string]string{
			"rel":   "search",
			"type":  "application/opensearchdescription+xml",
			"href":  d.Config.URL + o.URL,
			"title": o.Title,
		})
	}

	w.writeElement("title", d.Config.Title, nil)
	w.writeElement("subtitle", d.Config.Subtitle, nil)
	w.writeElement("update", time.Now().Format("2006-01-02T15:04:05Z07:00"), nil)

	addPostsToAtom(w, d)

	w.writeEndElement("feed")

	return w.bytes()
}

func addPostsToAtom(w *writer, d *data.Data) {
	for _, p := range d.Posts {
		w.writeStartElement("entry", nil, true)

		w.writeElement("id", p.Permalink, nil)

		w.writeCloseElement("link", map[string]string{
			"href": d.Config.URL + p.Permalink,
		})

		w.writeElement("title", p.Title, nil)

		t := time.Unix(p.Modified, 0)
		w.writeElement("update", t.Format("2006-01-02T15:04:05Z07:00"), nil)

		w.writeElement("summary", p.Summary, map[string]string{
			"type": "html",
		})

		w.writeEndElement("entry")
	}
}

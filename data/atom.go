// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/typing/data/xmlwriter"
)

// 用于生成一个符合 atom 规范的 XML 文本。
func (d *Data) buildAtom(conf *config) error {
	if conf.Atom == nil { // 不需要生成 atom
		return nil
	}

	w := xmlwriter.New()

	w.WriteStartElement("feed", map[string]string{
		"xmlns":            "http://www.w3.org/2005/Atom",
		"xmlns:opensearch": "http://a9.com/-/spec/opensearch/1.1/",
	})
	w.WriteElement("id", conf.URL, nil)
	w.WriteCloseElement("link", map[string]string{
		"href": conf.URL,
	})

	if conf.Opensearch != nil {
		o := conf.Opensearch
		w.WriteCloseElement("link", map[string]string{
			"rel":   "search",
			"type":  o.Type,
			"href":  d.url(o.URL),
			"title": o.Title,
		})
	}

	w.WriteElement("title", conf.Title, nil)
	w.WriteElement("subtitle", conf.Subtitle, nil)
	w.WriteElement("update", d.Created.Format(time.RFC3339), nil)

	addPostsToAtom(w, d)

	w.WriteEndElement("feed")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.Atom = &Feed{
		Title:   conf.Atom.Title,
		URL:     conf.Atom.URL,
		Type:    conf.Atom.Type,
		Content: bs,
	}

	return nil
}

func addPostsToAtom(w *xmlwriter.XMLWriter, d *Data) {
	for _, p := range d.Posts {
		w.WriteStartElement("entry", nil)

		w.WriteElement("id", p.Permalink, nil)

		w.WriteCloseElement("link", map[string]string{
			"href": d.url(p.Permalink),
		})

		w.WriteElement("title", p.Title, nil)

		w.WriteElement("update", p.Modified.Format(time.RFC3339), nil)

		w.WriteElement("summary", p.Summary, map[string]string{
			"type": "html",
		})

		w.WriteEndElement("entry")
	}
}

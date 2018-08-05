// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/data/xmlwriter"
	"github.com/issue9/web"
)

// 生成一个符合 RSS 规范的 XML 文本。
func (d *Data) buildRSS(conf *loader.Config) error {
	if conf.RSS == nil {
		return nil
	}

	w := xmlwriter.New()

	w.WriteStartElement("rss", map[string]string{
		"version":    "2.0",
		"xmlns:atom": "http://www.w3.org/2005/Atom",
	})
	w.WriteStartElement("channel", nil)

	w.WriteElement("title", conf.Title, nil)
	w.WriteElement("description", conf.Subtitle, nil)
	w.WriteElement("link", web.URL(""), nil)

	if conf.Opensearch != nil {
		w.WriteCloseElement("atom:link", map[string]string{
			"rel":   "search",
			"type":  conf.Opensearch.Type,
			"title": conf.Opensearch.Title,
			"href":  web.URL(conf.Opensearch.URL),
		})
	}

	addPostsToRSS(w, d)

	w.WriteEndElement("channel")
	w.WriteEndElement("rss")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.RSS = &Feed{
		Title:   conf.RSS.Title,
		URL:     conf.RSS.URL,
		Type:    conf.RSS.Type,
		Content: bs,
	}

	return nil
}

func addPostsToRSS(w *xmlwriter.XMLWriter, d *Data) {
	for _, p := range d.Posts {
		w.WriteStartElement("item", nil)

		w.WriteElement("link", web.URL(p.Permalink), nil)
		w.WriteElement("title", p.Title, nil)
		w.WriteElement("pubDate", p.Created.Format(time.RFC1123), nil)
		w.WriteElement("description", p.Summary, nil)

		w.WriteEndElement("item")
	}
}

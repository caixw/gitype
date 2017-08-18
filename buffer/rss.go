// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package buffer

import "time"

// 生成一个符合 rss 规范的 XML 文本。
func (buf *Buffer) buildRSS() error {
	conf := buf.Data.Config
	if conf.RSS == nil {
		return nil
	}

	w := newWrite()

	w.writeStartElement("rss", map[string]string{
		"version":    "2.0",
		"xmlns:atom": "http://www.w3.org/2005/Atom",
	})
	w.writeStartElement("channel", nil)

	w.writeElement("title", conf.Title, nil)
	w.writeElement("description", conf.Subtitle, nil)
	w.writeElement("link", conf.URL, nil)

	if conf.Opensearch != nil {
		w.writeCloseElement("atom:link", map[string]string{
			"rel":   "search",
			"type":  conf.Opensearch.Type,
			"title": conf.Opensearch.Title,
			"href":  conf.URL + conf.Opensearch.URL,
		})
	}

	addPostsToRSS(w, buf)

	w.writeEndElement("channel")
	w.writeEndElement("rss")

	bs, err := w.bytes()
	if err != nil {
		return err
	}
	buf.RSS = bs
	return nil
}

func addPostsToRSS(w *xmlWriter, buf *Buffer) {
	for _, p := range buf.Data.Posts {
		w.writeStartElement("item", nil)

		w.writeElement("link", buf.Data.Config.URL+p.Permalink, nil)
		w.writeElement("title", p.Title, nil)
		w.writeElement("pubDate", formatUnix(p.Created, time.RFC1123), nil)
		w.writeElement("description", p.Summary, nil)

		w.writeEndElement("item")
	}
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import "time"

// 用于生成一个符合 atom 规范的 XML 文本。
func (buf *Client) buildAtom() error {
	conf := buf.Data.Config
	if conf.Atom == nil { // 不需要生成 atom
		return nil
	}

	w := newWrite()

	w.writeStartElement("feed", map[string]string{
		"xmlns":            "http://www.w3.org/2005/Atom",
		"xmlns:opensearch": "http://a9.com/-/spec/opensearch/1.1/",
	})
	w.writeElement("id", conf.URL, nil)
	w.writeCloseElement("link", map[string]string{
		"href": conf.URL,
	})

	if conf.Opensearch != nil {
		o := conf.Opensearch
		w.writeCloseElement("link", map[string]string{
			"rel":   "search",
			"type":  o.Type,
			"href":  conf.URL + o.URL,
			"title": o.Title,
		})
	}

	w.writeElement("title", conf.Title, nil)
	w.writeElement("subtitle", conf.Subtitle, nil)
	w.writeElement("update", formatUnix(buf.Created, time.RFC3339), nil)

	addPostsToAtom(w, buf)

	w.writeEndElement("feed")

	bs, err := w.bytes()
	if err != nil {
		return err
	}
	buf.Atom = bs
	return nil
}

func addPostsToAtom(w *xmlWriter, buf *Client) {
	for _, p := range buf.Data.Posts {
		w.writeStartElement("entry", nil)

		w.writeElement("id", p.Permalink, nil)

		w.writeCloseElement("link", map[string]string{
			"href": buf.Data.Config.URL + p.Permalink,
		})

		w.writeElement("title", p.Title, nil)

		w.writeElement("update", formatUnix(p.Modified, time.RFC3339), nil)

		w.writeElement("summary", p.Summary, map[string]string{
			"type": "html",
		})

		w.writeEndElement("entry")
	}
}

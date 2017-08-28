// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"time"
)

func (client *Client) initAtom() error {
	if client.data.Config.Atom == nil { // 不需要生成 atom
		return nil
	}

	if err := client.buildAtom(); err != nil {
		return err
	}

	conf := client.data.Config
	client.patterns = append(client.patterns, conf.Atom.URL)
	client.mux.GetFunc(conf.Atom.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, conf.Atom.Type)
		w.Write(client.atom)
	}))
	return nil
}

// 用于生成一个符合 atom 规范的 XML 文本。
func (client *Client) buildAtom() error {
	conf := client.data.Config
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
			"href":  client.url(o.URL),
			"title": o.Title,
		})
	}

	w.writeElement("title", conf.Title, nil)
	w.writeElement("subtitle", conf.Subtitle, nil)
	w.writeElement("update", formatUnix(client.Created, time.RFC3339), nil)

	addPostsToAtom(w, client)

	w.writeEndElement("feed")

	bs, err := w.bytes()
	if err != nil {
		return err
	}
	client.atom = bs

	return nil
}

func addPostsToAtom(w *xmlWriter, client *Client) {
	for _, p := range client.data.Posts {
		w.writeStartElement("entry", nil)

		w.writeElement("id", p.Permalink, nil)

		w.writeCloseElement("link", map[string]string{
			"href": client.url(p.Permalink),
		})

		w.writeElement("title", p.Title, nil)

		w.writeElement("update", formatUnix(p.Modified, time.RFC3339), nil)

		w.writeElement("summary", p.Summary, map[string]string{
			"type": "html",
		})

		w.writeEndElement("entry")
	}
}

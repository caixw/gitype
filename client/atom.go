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
			"href":  client.url(o.URL),
			"title": o.Title,
		})
	}

	w.WriteElement("title", conf.Title, nil)
	w.WriteElement("subtitle", conf.Subtitle, nil)
	w.WriteElement("update", client.Created.Format(time.RFC3339), nil)

	addPostsToAtom(w, client)

	w.WriteEndElement("feed")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	client.atom = bs

	return nil
}

func addPostsToAtom(w *XMLWriter, client *Client) {
	for _, p := range client.data.Posts {
		w.WriteStartElement("entry", nil)

		w.WriteElement("id", p.Permalink, nil)

		w.WriteCloseElement("link", map[string]string{
			"href": client.url(p.Permalink),
		})

		w.WriteElement("title", p.Title, nil)

		w.WriteElement("update", p.Modified.Format(time.RFC3339), nil)

		w.WriteElement("summary", p.Summary, map[string]string{
			"type": "html",
		})

		w.WriteEndElement("entry")
	}
}

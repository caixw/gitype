// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"time"
)

func (client *Client) initRSS() error {
	conf := client.data.Config
	if conf.RSS == nil {
		return nil
	}

	if err := client.buildRSS(); err != nil {
		return err
	}

	client.patterns = append(client.patterns, conf.RSS.URL)
	client.mux.GetFunc(conf.RSS.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, conf.RSS.Type)
		w.Write(client.rss)
	}))

	return nil
}

// 生成一个符合 rss 规范的 XML 文本。
func (client *Client) buildRSS() error {
	conf := client.data.Config
	w := newWrite()

	w.WriteStartElement("rss", map[string]string{
		"version":    "2.0",
		"xmlns:atom": "http://www.w3.org/2005/Atom",
	})
	w.WriteStartElement("channel", nil)

	w.WriteElement("title", conf.Title, nil)
	w.WriteElement("description", conf.Subtitle, nil)
	w.WriteElement("link", conf.URL, nil)

	if conf.Opensearch != nil {
		w.WriteCloseElement("atom:link", map[string]string{
			"rel":   "search",
			"type":  conf.Opensearch.Type,
			"title": conf.Opensearch.Title,
			"href":  client.url(conf.Opensearch.URL),
		})
	}

	addPostsToRSS(w, client)

	w.WriteEndElement("channel")
	w.WriteEndElement("rss")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	client.rss = bs

	return nil
}

func addPostsToRSS(w *XMLWriter, client *Client) {
	for _, p := range client.data.Posts {
		w.WriteStartElement("item", nil)

		w.WriteElement("link", client.url(p.Permalink), nil)
		w.WriteElement("title", p.Title, nil)
		w.WriteElement("pubDate", p.Created.Format(time.RFC1123), nil)
		w.WriteElement("description", p.Summary, nil)

		w.WriteEndElement("item")
	}
}

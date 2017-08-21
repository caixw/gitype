// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"time"
)

func (client *Client) initRSS() error {
	conf := client.Data.Config
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
	conf := client.Data.Config
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

	addPostsToRSS(w, client)

	w.writeEndElement("channel")
	w.writeEndElement("rss")

	bs, err := w.bytes()
	if err != nil {
		return err
	}
	client.rss = bs

	return nil
}

func addPostsToRSS(w *xmlWriter, buf *Client) {
	for _, p := range buf.Data.Posts {
		w.writeStartElement("item", nil)

		w.writeElement("link", buf.Data.Config.URL+p.Permalink, nil)
		w.writeElement("title", p.Title, nil)
		w.writeElement("pubDate", formatUnix(p.Created, time.RFC1123), nil)
		w.writeElement("description", p.Summary, nil)

		w.writeEndElement("item")
	}
}

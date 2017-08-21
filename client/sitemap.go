// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/vars"
)

// 生成一个符合 sitemap 规范的 XML 文本。
func (client *Client) buildSitemap() error {
	conf := client.Data.Config
	if conf.Sitemap == nil {
		return nil
	}
	w := newWrite()

	if len(conf.Sitemap.XslURL) > 0 {
		w.writePI("xml-stylesheet", map[string]string{
			"type": "text/xsl",
			"href": conf.Sitemap.XslURL,
		})
	}

	w.writeStartElement("urlset", map[string]string{
		"xmlns": "http://www.sitemaps.org/schemas/sitemap/0.9",
	})

	addPostsToSitemap(w, client)

	if conf.Sitemap.EnableTag {
		addTagsToSitemap(w, client)
	}

	w.writeEndElement("urlset")

	bs, err := w.bytes()
	if err != nil {
		return err
	}
	client.sitemap = bs

	client.patterns = append(client.patterns, conf.Sitemap.URL)
	client.mux.GetFunc(conf.Sitemap.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
		setContentType(w, conf.Sitemap.Type)
		w.Write(client.sitemap)
	}))

	return nil
}

func addPostsToSitemap(w *xmlWriter, buf *Client) {
	sitemap := buf.Data.Config.Sitemap
	for _, p := range buf.Data.Posts {
		loc := buf.Data.Config.URL + p.Permalink
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *xmlWriter, buf *Client) error {
	sitemap := buf.Data.Config.Sitemap

	loc := buf.Data.Config.URL + vars.TagsURL()
	addItemToSitemap(w, loc, sitemap.TagChangefreq, buf.Created, sitemap.TagPriority)

	for _, tag := range buf.Data.Tags {
		loc = buf.Data.Config.URL + tag.Permalink
		addItemToSitemap(w, loc, sitemap.TagChangefreq, tag.Modified, sitemap.TagPriority)
	}
	return nil
}

func addItemToSitemap(w *xmlWriter, loc, changefreq string, lastmod int64, priority float64) {
	w.writeStartElement("url", nil)

	w.writeElement("loc", loc, nil)
	w.writeElement("lastmod", formatUnix(lastmod, time.RFC3339), nil)
	w.writeElement("changefreq", changefreq, nil)
	w.writeElement("priority", strconv.FormatFloat(priority, 'f', 1, 32), nil)

	w.writeEndElement("url")
}

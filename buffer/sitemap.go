// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package buffer

import (
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// 生成一个符合 sitemap 规范的 XML 文本。
func buildSitemap(d *data.Data) ([]byte, error) {
	w := newWrite()

	if len(d.Config.Sitemap.XslURL) > 0 {
		w.writePI("xml-stylesheet", map[string]string{
			"type": "text/xsl",
			"href": d.Config.Sitemap.XslURL,
		})
	}

	w.writeStartElement("urlset", map[string]string{
		"xmlns": "http://www.sitemaps.org/schemas/sitemap/0.9",
	})

	addPostsToSitemap(w, d)

	if d.Config.Sitemap.EnableTag {
		addTagsToSitemap(w, d)
	}

	w.writeEndElement("urlset")

	return w.bytes()
}

func addPostsToSitemap(w *xmlWriter, d *data.Data) {
	sitemap := d.Config.Sitemap
	for _, p := range d.Posts {
		loc := d.Config.URL + p.Permalink
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *xmlWriter, d *data.Data) error {
	sitemap := d.Config.Sitemap

	loc := d.Config.URL + vars.TagsURL()
	addItemToSitemap(w, loc, sitemap.TagChangefreq, time.Now().Unix(), sitemap.TagPriority)

	for _, tag := range d.Tags {
		loc = d.Config.URL + tag.Permalink
		addItemToSitemap(w, loc, sitemap.TagChangefreq, tag.Modified, sitemap.TagPriority)
	}
	return nil
}

func addItemToSitemap(w *xmlWriter, loc, changefreq string, lastmod int64, priority float64) {
	w.writeStartElement("url", nil)

	w.writeElement("loc", loc, nil)
	t := time.Unix(lastmod, 0)
	w.writeElement("lastmod", t.Format(time.RFC3339), nil)
	w.writeElement("changefreq", changefreq, nil)
	w.writeElement("priority", strconv.FormatFloat(priority, 'f', 1, 32), nil)

	w.writeEndElement("url")
}

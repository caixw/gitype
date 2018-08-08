// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/data/xmlwriter"
	"github.com/caixw/gitype/vars"
	"github.com/issue9/web"
)

// 生成一个符合 sitemap 规范的 XML 文本。
func (d *Data) buildSitemap(conf *loader.Config) error {
	if conf.Sitemap == nil {
		return nil
	}

	w := xmlwriter.New()

	if len(conf.Sitemap.XslURL) > 0 {
		w.WritePI("xml-stylesheet", map[string]string{
			"type": "text/xsl",
			"href": conf.Sitemap.XslURL,
		})
	}

	w.WriteStartElement("urlset", map[string]string{
		"xmlns": "http://www.sitemaps.org/schemas/sitemap/0.9",
	})

	addPostsToSitemap(w, d, conf)

	// archives.html
	loc := web.URL(vars.ArchivesURL())
	addItemToSitemap(w, loc, conf.Sitemap.Changefreq, d.Created, conf.Sitemap.Priority)

	// links.html
	loc = web.URL(vars.LinksURL())
	addItemToSitemap(w, loc, conf.Sitemap.Changefreq, d.Created, conf.Sitemap.Priority)

	if conf.Sitemap.EnableTag {
		addTagsToSitemap(w, d, conf)
	}

	w.WriteEndElement("urlset")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.Sitemap = &Feed{
		URL:     conf.Sitemap.URL,
		Type:    conf.Sitemap.Type,
		Content: bs,
	}

	return nil
}

func addPostsToSitemap(w *xmlwriter.XMLWriter, d *Data, conf *loader.Config) {
	sitemap := conf.Sitemap
	for _, p := range d.Posts {
		loc := web.URL(p.Permalink)
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *xmlwriter.XMLWriter, d *Data, conf *loader.Config) error {
	sitemap := conf.Sitemap

	loc := web.URL(vars.TagsURL())
	addItemToSitemap(w, loc, sitemap.Changefreq, d.Created, sitemap.Priority)

	for _, tag := range d.Tags {
		loc = web.URL(vars.TagURL(tag.Slug, 1))
		addItemToSitemap(w, loc, sitemap.Changefreq, tag.Modified, sitemap.Priority)
	}
	return nil
}

func addItemToSitemap(w *xmlwriter.XMLWriter, loc, changefreq string, lastmod time.Time, priority float64) {
	w.WriteStartElement("url", nil)

	w.WriteElement("loc", loc, nil)
	w.WriteElement("lastmod", lastmod.Format(time.RFC3339), nil)
	w.WriteElement("changefreq", changefreq, nil)
	w.WriteElement("priority", strconv.FormatFloat(priority, 'f', 1, 32), nil)

	w.WriteEndElement("url")
}

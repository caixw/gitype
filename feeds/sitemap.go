// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package feeds 包提供了根据数据生成 sitemap，atom，rss 的功能。
package feeds

import (
	"bytes"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
)

const (
	sitemapHeader = `<?xml version="1.0" encoding="utf-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	sitemapFooter = `</urlset>`
)

// BuildSitemap 生成一个符合 sitemap 规范的 XML 文本 buffer。
func BuildSitemap(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(sitemapHeader); err != nil {
		return nil, err
	}

	if err := addPostsToSitemap(buf, d); err != nil {
		return nil, err
	}

	if d.Config.Sitemap.EnableTag {
		if err := addTagsToSitemap(buf, d); err != nil {
			return nil, err
		}
	}

	if _, err := buf.WriteString(sitemapFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

func addPostsToSitemap(buf *bytes.Buffer, d *data.Data) error {
	sitemap := d.Config.Sitemap
	for _, p := range d.Posts {
		loc := d.Config.URL + p.Permalink
		addItemToSitemap(buf, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
	return nil
}

func addTagsToSitemap(buf *bytes.Buffer, d *data.Data) error {
	now := time.Now().Unix()
	sitemap := d.Config.Sitemap

	loc := d.Config.URL + d.Config.URLS.Tags + d.Config.URLS.Suffix
	addItemToSitemap(buf, loc, sitemap.TagChangefreq, now, sitemap.TagPriority)

	for _, tag := range d.Tags {
		loc = d.Config.URL + tag.Permalink
		addItemToSitemap(buf, loc, sitemap.TagChangefreq, now, sitemap.TagPriority)
	}
	return nil
}

func addItemToSitemap(buf *bytes.Buffer, loc, changefreq string, lastmod int64, priority float64) {
	buf.WriteString("<url>\n")

	buf.WriteString("<loc>")
	buf.WriteString(loc)
	buf.WriteString("</loc>\n")

	t := time.Unix(lastmod, 0)
	buf.WriteString("<lastmod>")
	buf.WriteString(t.Format("2006-01-02T15:04:05-07:00"))
	buf.WriteString("</lastmod>\n")

	buf.WriteString("<changefreq>")
	buf.WriteString(changefreq)
	buf.WriteString("</changefreq>\n")

	buf.WriteString("<priority>")
	buf.WriteString(strconv.FormatFloat(priority, 'f', 1, 32))
	buf.WriteString("</priority>\n")

	buf.WriteString("</url>\n")
}

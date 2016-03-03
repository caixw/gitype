// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feeds

import (
	"bytes"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
)

const (
	sitemapHeader = `<?xml version="1.0" encoding="utf-8"?>
<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>`

	sitemapFooter = `</urlset>`
)

// Build 构建一个sitemap.xml文件到sitemapPath文件中，若该文件已经存在，则覆盖。
func BuildSitemap(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(sitemapHeader); err != nil {
		return nil, err
	}

	if err := addPostsToSitemap(buf, d); err != nil {
		return nil, err
	}

	if err := addTagsToSitemap(buf, d); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString(sitemapFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

func addPostsToSitemap(buf *bytes.Buffer, d *data.Data) error {
	sitemap := d.Config.Sitemap
	for _, p := range d.Posts {
		addItemToSitemap(buf, p.Permalink, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
	return nil
}

func addTagsToSitemap(buf *bytes.Buffer, d *data.Data) error {
	now := time.Now().Unix()
	sitemap := d.Config.Sitemap
	for _, tag := range d.Tags {
		addItemToSitemap(buf, tag.Permalink, sitemap.TagChangefreq, now, sitemap.TagPriority)
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
	buf.WriteString(t.Format("2006-01-02T15:04:05+08:00"))
	buf.WriteString("</lastmod>\n")

	buf.WriteString("<changefreq>")
	buf.WriteString(changefreq)
	buf.WriteString("</changefreq>\n")

	buf.WriteString("<priority>")
	buf.WriteString(strconv.FormatFloat(priority, 'f', 1, 32))
	buf.WriteString("</priority>\n")

	buf.WriteString("</url>\n")
}

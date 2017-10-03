// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/url"
)

const contentTypeXML = "application/xml"

type sitemapConfig struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type,omitempty"`

	XslURL     string  `yaml:"xslURL,omitempty"`    // 为 sitemap 指定一个 xsl 文件
	Priority   float64 `yaml:"priority"`            // 默认的优先级
	Changefreq string  `yaml:"changefreq"`          // 默认的更新频率
	EnableTag  bool    `yaml:"enableTag,omitempty"` // 是否将标签相关的页面写入 sitemap

	// 文章可以指定一个专门的值
	PostPriority   float64 `yaml:"postPriority"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

// 生成一个符合 sitemap 规范的 XML 文本。
func (d *Data) buildSitemap(conf *config) error {
	if conf.Sitemap == nil {
		return nil
	}

	w := helper.NewWriter()

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
	loc := d.URL(url.Archives())
	addItemToSitemap(w, loc, conf.Sitemap.Changefreq, d.Created, conf.Sitemap.Priority)

	// links.html
	loc = d.URL(url.Links())
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

func addPostsToSitemap(w *helper.XMLWriter, d *Data, conf *config) {
	sitemap := conf.Sitemap
	for _, p := range d.Posts {
		loc := d.URL(p.Permalink)
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *helper.XMLWriter, d *Data, conf *config) error {
	sitemap := conf.Sitemap

	loc := d.URL(url.Tags())
	addItemToSitemap(w, loc, sitemap.Changefreq, d.Created, sitemap.Priority)

	for _, tag := range d.Tags {
		loc = d.URL(url.Tag(tag.Slug, 1))
		addItemToSitemap(w, loc, sitemap.Changefreq, tag.Modified, sitemap.Priority)
	}
	return nil
}

func addItemToSitemap(w *helper.XMLWriter, loc, changefreq string, lastmod time.Time, priority float64) {
	w.WriteStartElement("url", nil)

	w.WriteElement("loc", loc, nil)
	w.WriteElement("lastmod", lastmod.Format(time.RFC3339), nil)
	w.WriteElement("changefreq", changefreq, nil)
	w.WriteElement("priority", strconv.FormatFloat(priority, 'f', 1, 32), nil)

	w.WriteEndElement("url")
}

// 检测 sitemap 取值是否正确
func (s *sitemapConfig) sanitize() *helper.FieldError {
	switch {
	case len(s.URL) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "sitemap.url"}
	case s.Priority > 1 || s.Priority < 0:
		return &helper.FieldError{Message: "介于[0,1]之间的浮点数", Field: "sitemap.priority"}
	case s.PostPriority > 1 || s.PostPriority < 0:
		return &helper.FieldError{Message: "介于[0,1]之间的浮点数", Field: "sitemap.postPriority"}
	case !isChangereq(s.Changefreq):
		return &helper.FieldError{Message: "取值不正确", Field: "sitemap.changefreq"}
	case !isChangereq(s.PostChangefreq):
		return &helper.FieldError{Message: "取值不正确", Field: "sitemap.postChangefreq"}
	}

	if len(s.Type) == 0 {
		s.Type = contentTypeXML
	}

	return nil
}

var changereqs = []string{
	"never",
	"yearly",
	"monthly",
	"weekly",
	"daily",
	"hourly",
	"always",
}

func isChangereq(val string) bool {
	for _, v := range changereqs {
		if v == val {
			return true
		}
	}
	return false
}

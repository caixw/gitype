// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"
	"time"

	"github.com/caixw/typing/data/xmlwriter"
	"github.com/caixw/typing/vars"
)

// Sitemap 的相关参数
type Sitemap struct {
	URL     string // 展示给用户的地址，不能包含域名
	Type    string // mime type
	Content []byte // 实际内容
}

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
	loc := d.url(vars.ArchivesURL())
	addItemToSitemap(w, loc, conf.Sitemap.Changefreq, d.Created, conf.Sitemap.Priority)

	if conf.Sitemap.EnableTag {
		addTagsToSitemap(w, d, conf)
	}

	w.WriteEndElement("urlset")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.Sitemap = &Sitemap{
		URL:     conf.Sitemap.URL,
		Type:    conf.Sitemap.Type,
		Content: bs,
	}

	return nil
}

func addPostsToSitemap(w *xmlwriter.XMLWriter, d *Data, conf *config) {
	sitemap := conf.Sitemap
	for _, p := range d.Posts {
		loc := d.url(p.Permalink)
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *xmlwriter.XMLWriter, d *Data, conf *config) error {
	sitemap := conf.Sitemap

	loc := d.url(vars.TagsURL())
	addItemToSitemap(w, loc, sitemap.Changefreq, d.Created, sitemap.Priority)

	for _, tag := range d.Tags {
		loc = d.url(vars.TagURL(tag.Slug, 1))
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

// 检测 sitemap 取值是否正确
func (s *sitemapConfig) sanitize() *FieldError {
	switch {
	case len(s.URL) == 0:
		return &FieldError{Message: "不能为空", Field: "Sitemap.URL"}
	case s.Priority > 1 || s.Priority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "Sitemap.priority"}
	case s.PostPriority > 1 || s.PostPriority < 0:
		return &FieldError{Message: "介于[0,1]之间的浮点数", Field: "Sitemap.PostPriority"}
	case !isChangereq(s.Changefreq):
		return &FieldError{Message: "取值不正确", Field: "Sitemap.changefreq"}
	case !isChangereq(s.PostChangefreq):
		return &FieldError{Message: "取值不正确", Field: "Sitemap.PostChangefreq"}
	}

	if len(s.Type) == 0 {
		s.Type = vars.ContentTypeXML
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

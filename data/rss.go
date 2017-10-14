// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/gitype/helper"
)

const (
	contentTypeAtom = "application/atom+xml"
	contentTypeRSS  = "application/rss+xml"
)

// RSS 和 Atom 相关的配置项
type rssConfig struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Type  string `yaml:"type,omitempty"`
	Size  int    `yaml:"size"` // 显示数量
}

// 生成一个符合 RSS 规范的 XML 文本。
func (d *Data) buildRSS(conf *config) error {
	if conf.RSS == nil {
		return nil
	}

	w := helper.NewWriter()

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
			"href":  d.BuildURL(conf.Opensearch.URL),
		})
	}

	addPostsToRSS(w, d)

	w.WriteEndElement("channel")
	w.WriteEndElement("rss")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.RSS = &Feed{
		Title:   conf.RSS.Title,
		URL:     conf.RSS.URL,
		Type:    conf.RSS.Type,
		Content: bs,
	}

	return nil
}

func addPostsToRSS(w *helper.XMLWriter, d *Data) {
	for _, p := range d.Posts {
		w.WriteStartElement("item", nil)

		w.WriteElement("link", d.BuildURL(p.Permalink), nil)
		w.WriteElement("title", p.Title, nil)
		w.WriteElement("pubDate", p.Created.Format(time.RFC1123), nil)
		w.WriteElement("description", p.Summary, nil)

		w.WriteEndElement("item")
	}
}

func (rss *rssConfig) sanitize(conf *config, typ string) *helper.FieldError {
	if rss.Size <= 0 {
		return &helper.FieldError{Message: "必须大于 0", Field: typ + ".Size"}
	}
	if len(rss.URL) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: typ + ".URL"}
	}

	switch typ {
	case "rss":
		rss.Type = contentTypeRSS
	case "atom":
		rss.Type = contentTypeAtom
	default:
		panic("无效的 typ 值")
	}

	if len(rss.Title) == 0 {
		rss.Title = conf.Title
	}

	return nil
}

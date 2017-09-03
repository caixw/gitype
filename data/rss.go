// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"time"

	"github.com/caixw/typing/data/xmlwriter"
	"github.com/caixw/typing/vars"
)

// RSS 的相关配置
type RSS struct {
	Title   string
	URL     string
	Type    string
	Content []byte
}

type rssConfig struct {
	Title string `yaml:"title"`          // 标题
	URL   string `yaml:"url"`            // 地址，不能包含域名
	Type  string `yaml:"type,omitempty"` // mimeType
	Size  int    `yaml:"size"`           // 显示数量
}

// 生成一个符合 rss 规范的 XML 文本。
func (d *Data) buildRSS() error {
	conf := d.Config
	w := xmlwriter.New()

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
			"href":  d.url(conf.Opensearch.URL),
		})
	}

	addPostsToRSS(w, d)

	w.WriteEndElement("channel")
	w.WriteEndElement("rss")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.RSS = &RSS{
		Title:   d.Config.RSS.Title,
		URL:     d.Config.RSS.URL,
		Type:    d.Config.RSS.Type,
		Content: bs,
	}

	return nil
}

func addPostsToRSS(w *xmlwriter.XMLWriter, d *Data) {
	for _, p := range d.Posts {
		w.WriteStartElement("item", nil)

		w.WriteElement("link", d.url(p.Permalink), nil)
		w.WriteElement("title", p.Title, nil)
		w.WriteElement("pubDate", p.Created.Format(time.RFC1123), nil)
		w.WriteElement("description", p.Summary, nil)

		w.WriteEndElement("item")
	}
}

func (rss *rssConfig) sanitize(conf *Config, typ string) *FieldError {
	if rss.Size <= 0 {
		return &FieldError{Message: "必须大于 0", Field: typ + ".Size"}
	}
	if len(rss.URL) == 0 {
		return &FieldError{Message: "不能为空", Field: typ + ".URL"}
	}

	switch typ {
	case "rss":
		rss.Type = vars.ContentTypeRSS
	case "atom":
		rss.Type = vars.ContentTypeAtom
	default:
		panic("无效的 typ 值")
	}

	if len(rss.Title) == 0 {
		rss.Title = conf.Title
	}

	return nil
}

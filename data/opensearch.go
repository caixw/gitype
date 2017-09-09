// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"net/http"

	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
)

type opensearchConfig struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type,omitempty"`
	Title string `yaml:"title,omitempty"`

	ShortName   string `yaml:"shortName"`
	Description string `yaml:"description"`
	LongName    string `yaml:"longName,omitempty"`
	Image       *Icon  `yaml:"image,omitempty"`
}

// 用于生成一个符合 atom 规范的 XML 文本。
func (d *Data) buildOpensearch(conf *config) error {
	if conf.Opensearch == nil {
		return nil
	}

	w := helper.NewWriter()
	o := conf.Opensearch

	w.WriteStartElement("OpenSearchDescription", map[string]string{
		"xmlns": "http://a9.com/-/spec/opensearch/1.1/",
	})

	w.WriteElement("InputEncoding", "UTF-8", nil)
	w.WriteElement("OutputEncoding", "UTF-8", nil)
	w.WriteElement("ShortName", o.ShortName, nil)
	w.WriteElement("Description", o.Description, nil)

	if len(o.LongName) > 0 {
		w.WriteElement("LongName", o.LongName, nil)
	}

	if o.Image != nil {
		w.WriteElement("Image", o.Image.URL, map[string]string{
			"type": o.Image.Type,
		})
	}

	w.WriteCloseElement("Url", map[string]string{
		"type":   conf.Type,
		"method": http.MethodGet,
		// 需要全链接，否则 Firefox 的搜索框不认。
		// https://github.com/caixw/typing/issues/18
		"template": d.URL(vars.SearchURL("{searchTerms}", 0)),
	})

	w.WriteElement("Developer", vars.Name, nil)
	w.WriteElement("Language", conf.Language, nil)

	w.WriteEndElement("OpenSearchDescription")

	bs, err := w.Bytes()
	if err != nil {
		return err
	}
	d.Opensearch = &Feed{
		URL:     o.URL,
		Type:    o.Type,
		Title:   o.Title,
		Content: bs,
	}

	return nil
}

// 检测 opensearch 取值是否正确
func (s *opensearchConfig) sanitize(conf *config) *helper.FieldError {
	switch {
	case len(s.URL) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.url"}
	case len(s.ShortName) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.shortName"}
	case len(s.Description) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.description"}
	}

	if len(s.Type) == 0 {
		s.Type = vars.ContentTypeOpensearch
	}

	if s.Image == nil && conf.Icon != nil {
		s.Image = conf.Icon
	}

	return nil
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"net/http"

	"github.com/issue9/web"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/vars"
)

// 用于生成一个符合 opensearch 规范的 XML 文本。
func (d *Data) buildOpensearch(conf *loader.Config) error {
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
		// https://github.com/caixw/gitype/issues/18
		"template": web.URL(vars.SearchURL("{searchTerms}", 0)),
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

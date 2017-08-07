// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// BuildOpensearch 用于生成一个符合 atom 规范的 XML 文本。
func BuildOpensearch(d *data.Data) ([]byte, error) {
	w := newWrite()

	w.writeStartElement("OpenSearchDescription", map[string]string{
		"xmlns": "http://a9.com/-/spec/opensearch/1.1/",
	})

	o := d.Config.Opensearch

	w.writeElement("InputEncoding", "UTF-8", nil)
	w.writeElement("OutputEncoding", "UTF-8", nil)
	w.writeElement("ShortName", o.ShortName, nil)
	w.writeElement("Description", o.Description, nil)

	if len(o.LongName) > 0 {
		w.writeElement("LongName", o.LongName, nil)
	}

	if o.Image != nil {
		w.writeElement("Image", o.Image.URL, map[string]string{
			"type": o.Image.Type,
		})
	}

	w.writeCloseElement("Url", map[string]string{
		"type":     "text/html",
		"template": vars.SearchURL("{searchTerms}", 0),
	})

	w.writeElement("Developer", vars.AppName, nil)
	w.writeElement("Language", d.Config.Language, nil)

	w.writeEndElement("OpenSearchDescription")

	return w.bytes()
}

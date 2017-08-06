// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"

	"mime"
	"path/filepath"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

const (
	opensearchHeader = `<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">`

	opensearchFooter = `</OpenSearchDescription>`
)

// BuildOpensearch 用于生成一个符合 atom 规范的 XML 文本 buffer。
func BuildOpensearch(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := &writer{
		buf: buf,
	}

	w.writeString(xmlHeader)
	w.writeString(opensearchHeader)

	o := d.Config.Opensearch

	w.writeElement("InputEncoding", "UTF-8", nil)
	w.writeElement("OutputEncoding", "UTF-8", nil)
	w.writeElement("ShortName", o.ShortName, nil)
	w.writeElement("Description", o.Description, nil)

	if len(o.LongName) > 0 {
		w.writeElement("LongName", o.LongName, nil)
	}

	if len(o.Image) > 0 {
		w.writeElement("Image", o.Image, map[string]string{
			"type": mime.TypeByExtension(filepath.Ext(o.Image)),
		})
	}

	w.writeString(`<Url type="text/html" template="`)
	w.writeString(vars.SearchURL("{searchTerms}", 0))
	w.writeString(`" />`)

	w.writeElement("Developer", vars.AppName, nil)
	w.writeElement("Language", d.Config.Language, nil)

	w.writeString(opensearchFooter)

	return buf, nil
}

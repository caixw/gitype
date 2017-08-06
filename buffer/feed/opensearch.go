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
	w := &errWriter{
		buf: buf,
	}

	w.writeString(xmlHeader)

	w.writeString(opensearchHeader)

	o := d.Config.Opensearch

	w.writeString("\n<InputEncoding>")
	w.writeString("UTF-8")
	w.writeString("</InputEncoding>\n")

	w.writeString("<OutputEncoding>")
	w.writeString("UTF-8")
	w.writeString("</OutputEncoding>\n")

	w.writeString("<ShortName>")
	w.writeString(o.ShortName)
	w.writeString("</ShortName>\n")

	w.writeString("<Description>")
	w.writeString(o.Description)
	w.writeString("</Description>\n")

	if len(o.LongName) > 0 {
		w.writeString("<LongName>")
		w.writeString(o.LongName)
		w.writeString("</LongName>\n")
	}

	if len(o.Image) > 0 {
		w.writeString(`<Image type="`)
		w.writeString(mime.TypeByExtension(filepath.Ext(o.Image)))
		w.writeString(`">`)
		w.writeString(o.Image)
		w.writeString("</Image>\n")
	}

	w.writeString(`<Url type="text/html" template="`)
	w.writeString(vars.SearchURL("{searchTerms}", 0))
	w.writeString(`" />`)

	w.writeString("<Developer>")
	w.writeString(vars.AppName)
	w.writeString("</Developer>\n")

	w.writeString("<Language>")
	w.writeString(d.Config.Language)
	w.writeString("</Language>\n")

	w.writeString(opensearchFooter)

	return buf, nil
}

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
	opensearchHeader = `<?xml version="1.0" encoding="utf-8"?>
<OpenSearchDescription xmlns="http://a9.com/-/spec/opensearch/1.1/">`

	opensearchFooter = `</OpenSearchDescription>`
)

// BuildOpensearch 用于生成一个符合 atom 规范的 XML 文本 buffer。
func BuildOpensearch(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(opensearchHeader); err != nil {
		return nil, err
	}

	o := d.Config.Opensearch

	buf.WriteString("\n<InputEncoding>")
	buf.WriteString("UTF-8")
	buf.WriteString("</InputEncoding>\n")

	buf.WriteString("<OutputEncoding>")
	buf.WriteString("UTF-8")
	buf.WriteString("</OutputEncoding>\n")

	buf.WriteString("<ShortName>")
	buf.WriteString(o.ShortName)
	buf.WriteString("</ShortName>\n")

	buf.WriteString("<Description>")
	buf.WriteString(o.Description)
	buf.WriteString("</Description>\n")

	if len(o.LongName) > 0 {
		buf.WriteString("<LongName>")
		buf.WriteString(o.LongName)
		buf.WriteString("</LongName>\n")
	}

	if len(o.Image) > 0 {
		buf.WriteString(`<Image type="`)
		buf.WriteString(mime.TypeByExtension(filepath.Ext(o.Image)))
		buf.WriteString(`">`)
		buf.WriteString(o.Image)
		buf.WriteString("</Image>\n")
	}

	buf.WriteString(`<Url type="text/html" template="`)
	buf.WriteString(vars.SearchURL("{searchTerms}", 0))
	buf.WriteString(`" />`)

	buf.WriteString("<Developer>")
	buf.WriteString(vars.AppName)
	buf.WriteString("</Developer>\n")

	buf.WriteString("<Language>")
	buf.WriteString(d.Config.Language)
	buf.WriteString("</Language>\n")

	if _, err := buf.WriteString(opensearchFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

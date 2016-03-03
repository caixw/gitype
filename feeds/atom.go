// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feeds

import (
	"bytes"
	"time"

	"github.com/caixw/typing/data"
)

const (
	atomHeader = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">`

	atomFooter = `</feed>`
)

// Build 构建一个atom.xml文件到atomPath文件中，若该文件已经存在，则覆盖。
func BuildAtom(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(atomHeader); err != nil {
		return nil, err
	}

	buf.WriteString("<id>")
	buf.WriteString(d.Config.URL)
	buf.WriteString("</id>\n")

	buf.WriteString("<link>")
	buf.WriteString(d.Config.URL)
	buf.WriteString("</link>\n")

	buf.WriteString("<title>")
	buf.WriteString(d.Config.Title)
	buf.WriteString("</title>\n")

	buf.WriteString("<subtitle>")
	buf.WriteString(d.Config.Subtitle)
	buf.WriteString("</subtitle>\n")

	buf.WriteString("<update>")
	buf.WriteString(time.Now().Format("2006-01-02T15:04:05Z07:00"))
	buf.WriteString("</update>\n")

	addPostsToAtom(buf, d)

	if _, err := buf.WriteString(atomFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

func addPostsToAtom(buf *bytes.Buffer, d *data.Data) {
	for _, p := range d.Posts {
		buf.WriteString("<entry>\n")

		buf.WriteString("<id>")
		buf.WriteString(p.Permalink)
		buf.WriteString("</id>\n")

		buf.WriteString("<link>")
		buf.WriteString(p.Permalink)
		buf.WriteString("</link>\n")

		buf.WriteString("<title>")
		buf.WriteString(p.Title)
		buf.WriteString("</title>\n")

		t := time.Unix(p.Modified, 0)
		buf.WriteString("<update>")
		buf.WriteString(t.Format("2006-01-02T15:04:05Z07:00"))
		buf.WriteString("</update>\n")

		buf.WriteString("<summary>")
		buf.WriteString(p.Summary)
		buf.WriteString("</summary>\n")

		buf.WriteString("</entry>\n")
	}
}

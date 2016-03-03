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
	rssHeader = `<?xml version="1.0" encoding="utf-8"?>
<rss version="2.0">
<channel>`

	rssFooter = `</channel>
</rss>`
)

// Build 构建一个rss.xml文件到rssPath文件中，若该文件已经存在，则覆盖。
func BuildRSS(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.WriteString(rssHeader); err != nil {
		return nil, err
	}

	addPostsToRSS(buf, d)

	if _, err := buf.WriteString(rssFooter); err != nil {
		return nil, err
	}

	return buf, nil
}

func addPostsToRSS(buf *bytes.Buffer, d *data.Data) {
	for _, p := range d.Posts {
		buf.WriteString("<item>\n")

		buf.WriteString("<link>")
		buf.WriteString(p.Permalink)
		buf.WriteString("</link>\n")

		buf.WriteString("<title>")
		buf.WriteString(p.Title)
		buf.WriteString("</title>\n")

		t := time.Unix(p.Created, 0)
		buf.WriteString("<pubDate>")
		buf.WriteString(t.Format(time.RFC1123))
		buf.WriteString("</pubDate>\n")

		buf.WriteString("<description>")
		buf.WriteString(p.Summary)
		buf.WriteString("</description>\n")

		buf.WriteString("</item>\n")
	}
}

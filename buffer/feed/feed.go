// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
)

const xmlHeader = `<?xml version="1.0" encoding="utf-8" ?>`

func writePI(buf *bytes.Buffer, name string, kv map[string]string) error {
	buf.WriteString("<?")
	buf.WriteString(name)

	for k, v := range kv {
		buf.WriteByte(' ')
		buf.WriteString(k)
		buf.WriteString(`="`)
		buf.WriteString(v)
		buf.WriteString(`"`)
	}

	buf.WriteString(" ?>")

	return nil
}

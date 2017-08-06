// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
)

const xmlHeader = `<?xml version="1.0" encoding="utf-8" ?>`

type errWriter struct {
	err error
	buf *bytes.Buffer
}

func (e *errWriter) writeString(str string) {
	if e.err != nil {
		return
	}

	_, e.err = e.buf.WriteString(str)
}

func (e *errWriter) writeByte(b byte) {
	if e.err != nil {
		return
	}

	e.err = e.buf.WriteByte(b)
}

func writePI(buf *bytes.Buffer, name string, kv map[string]string) error {
	w := &errWriter{
		buf: buf,
	}

	w.writeString("<?")
	w.writeString(name)

	for k, v := range kv {
		w.writeByte(' ')
		w.writeString(k)
		w.writeString(`="`)
		w.writeString(v)
		w.writeString(`"`)
	}

	w.writeString(" ?>")

	return w.err
}

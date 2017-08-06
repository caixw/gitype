// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
)

const xmlHeader = `<?xml version="1.0" encoding="utf-8" ?>`

type writer struct {
	err error
	buf *bytes.Buffer
}

func newWrite() *writer {
	w := &writer{
		buf: new(bytes.Buffer),
	}

	w.writeString(xmlHeader)

	return w
}

func (w *writer) writeString(str string) {
	if w.err != nil {
		return
	}

	_, w.err = w.buf.WriteString(str)
}

func (w *writer) writeByte(b byte) {
	if w.err != nil {
		return
	}

	w.err = w.buf.WriteByte(b)
}

func (w *writer) writeCloseElement(name string, attr map[string]string) {
	w.writeByte('<')
	w.writeString(name)
	for k, v := range attr {
		w.writeByte(' ')
		w.writeString(k)
		w.writeString(`="`)
		w.writeString(v)
		w.writeString(`"`)
	}
	w.writeString(" />")
}

func (w *writer) writeElement(name, val string, attr map[string]string) {
	w.writeByte('<')
	w.writeString(name)
	for k, v := range attr {
		w.writeByte(' ')
		w.writeString(k)
		w.writeString(`="`)
		w.writeString(v)
		w.writeString(`"`)
	}
	w.writeByte('>')

	w.writeString(val)

	w.writeString("</")
	w.writeString(name)
	w.writeString(">\n")
}

func (w *writer) writePI(name string, kv map[string]string) {
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
}

// 将内容转换成 []byte 并返回
func (w *writer) bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	return w.buf.Bytes(), nil
}

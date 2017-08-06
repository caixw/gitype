// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
)

// xml 操作类，简单地封装 bytes.Buffer。
type writer struct {
	err error
	buf *bytes.Buffer
}

func newWrite() *writer {
	w := &writer{
		buf: new(bytes.Buffer),
	}

	w.writePI("xml", map[string]string{
		"version":  "1.0",
		"encoding": "utf-8",
	})

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

func (w *writer) writeNewline() {
	w.writeByte('\n')
}

func (w *writer) writeStartElement(name string, attr map[string]string, newline bool) {
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

	if newline {
		w.writeNewline()
	}
}

func (w *writer) writeEndElement(name string) {
	w.writeString("</")
	w.writeString(name)
	w.writeByte('>')
	w.writeNewline()
}

// 写入一个自闭合的元素
// name 元素标签名；
// attr 元素的属性。
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
	w.writeNewline()
}

// 写入一个元素。
// name 元素标签名；
// val 元素内容；
// attr 元素的属性。
func (w *writer) writeElement(name, val string, attr map[string]string) {
	w.writeStartElement(name, attr, false)
	w.writeString(val)
	w.writeEndElement(name)
}

// 写入一个 PI 指令
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
	w.writeNewline()
}

// 将内容转换成 []byte 并返回
func (w *writer) bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	return w.buf.Bytes(), nil
}

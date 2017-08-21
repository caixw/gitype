// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"strings"
	"time"
)

// XML 要求 version 属于必须在其它属性之前
const xmlPI = `<?xml version="1.0" encoding="utf-8"?>`

func formatUnix(unix int64, format string) string {
	return time.Unix(unix, 0).Format(format)
}

// xml 操作类，简单地封装 bytes.Buffer。
type xmlWriter struct {
	buf    *bytes.Buffer
	err    error // 缓存 buf.Write* 系列函数的错误信息，并阻止其再次执行
	indent int   // 保存当前的缩进量
}

func newWrite() *xmlWriter {
	w := &xmlWriter{
		buf: bytes.NewBufferString(xmlPI),
	}

	w.writeByte('\n')

	return w
}

func (w *xmlWriter) writeString(str string) {
	if w.err != nil {
		return
	}

	_, w.err = w.buf.WriteString(str)
}

func (w *xmlWriter) writeByte(b byte) {
	if w.err != nil {
		return
	}

	w.err = w.buf.WriteByte(b)
}

func (w *xmlWriter) writeStartElement(name string, attr map[string]string) {
	w.startElement(name, attr, true)
}

// newline 是否换行
func (w *xmlWriter) startElement(name string, attr map[string]string, newline bool) {
	w.writeString(strings.Repeat(" ", w.indent*4))
	w.indent++

	w.writeByte('<')
	w.writeString(name)
	w.writeAttr(attr)
	w.writeByte('>')

	if newline {
		w.writeByte('\n')
	}
}

func (w *xmlWriter) writeEndElement(name string) {
	w.endElement(name, true)
}

// indent 是否需要填上缩时的字符
func (w *xmlWriter) endElement(name string, indent bool) {
	w.indent--
	if indent {
		w.writeString(strings.Repeat(" ", w.indent*4))
	}

	w.writeString("</")
	w.writeString(name)
	w.writeByte('>')

	w.writeByte('\n')
}

// 写入一个自闭合的元素
// name 元素标签名；
// attr 元素的属性。
func (w *xmlWriter) writeCloseElement(name string, attr map[string]string) {
	w.writeString(strings.Repeat(" ", w.indent*4))

	w.writeByte('<')
	w.writeString(name)
	w.writeAttr(attr)
	w.writeString(" />")

	w.writeByte('\n')
}

// 写入一个元素。
// name 元素标签名；
// val 元素内容；
// attr 元素的属性。
func (w *xmlWriter) writeElement(name, val string, attr map[string]string) {
	w.startElement(name, attr, false)
	w.writeString(val)
	w.endElement(name, false)
}

// 写入一个 PI 指令
func (w *xmlWriter) writePI(name string, kv map[string]string) {
	w.writeString("<?")
	w.writeString(name)
	w.writeAttr(kv)
	w.writeString("?>")

	w.writeByte('\n')
}

func (w *xmlWriter) writeAttr(attr map[string]string) {
	for k, v := range attr {
		w.writeByte(' ')
		w.writeString(k)
		w.writeString(`="`)
		w.writeString(v)
		w.writeByte('"')
	}
}

// 将内容转换成 []byte 并返回
func (w *xmlWriter) bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	return w.buf.Bytes(), nil
}

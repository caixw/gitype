// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xmlwriter

import (
	"bytes"
	"strings"
)

// XML 要求 version 属于必须在其它属性之前
const xmlPI = `<?xml version="1.0" encoding="utf-8"?>`

// XMLWriter 操作类，简单地封装 bytes.Buffer。
type XMLWriter struct {
	buf    *bytes.Buffer
	err    error // 缓存 buf.Write* 系列函数的错误信息，并阻止其再次执行
	indent int   // 保存当前的缩进量
}

// New 声明一个新的 XMLWriter
func New() *XMLWriter {
	w := &XMLWriter{
		buf: bytes.NewBufferString(xmlPI),
	}

	w.writeByte('\n')

	return w
}

func (w *XMLWriter) writeString(str string) {
	if w.err != nil {
		return
	}

	_, w.err = w.buf.WriteString(str)
}

func (w *XMLWriter) writeByte(b byte) {
	if w.err != nil {
		return
	}

	w.err = w.buf.WriteByte(b)
}

// WriteStartElement 写入一个开始元素
func (w *XMLWriter) WriteStartElement(name string, attr map[string]string) {
	w.startElement(name, attr, true)
}

// newline 是否换行
func (w *XMLWriter) startElement(name string, attr map[string]string, newline bool) {
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

// WriteEndElement 写入一个结束元素
func (w *XMLWriter) WriteEndElement(name string) {
	w.endElement(name, true)
}

// indent 是否需要填上缩进时的字符，如果不换行输出结束符，则不能输出缩进字符串
func (w *XMLWriter) endElement(name string, indent bool) {
	w.indent--
	if indent {
		w.writeString(strings.Repeat(" ", w.indent*4))
	}

	w.writeString("</")
	w.writeString(name)
	w.writeByte('>')

	w.writeByte('\n')
}

// WriteCloseElement 写入一个自闭合的元素
// name 元素标签名；
// attr 元素的属性。
func (w *XMLWriter) WriteCloseElement(name string, attr map[string]string) {
	w.writeString(strings.Repeat(" ", w.indent*4))

	w.writeByte('<')
	w.writeString(name)
	w.writeAttr(attr)
	w.writeString(" />")

	w.writeByte('\n')
}

// WriteElement 写入一个完整的元素。
// name 元素标签名；
// val 元素内容；
// attr 元素的属性。
func (w *XMLWriter) WriteElement(name, val string, attr map[string]string) {
	w.startElement(name, attr, false)
	w.writeString(val)
	w.endElement(name, false)
}

// WritePI 写入一个 PI 指令
func (w *XMLWriter) WritePI(name string, kv map[string]string) {
	w.writeString("<?")
	w.writeString(name)
	w.writeAttr(kv)
	w.writeString("?>")

	w.writeByte('\n')
}

func (w *XMLWriter) writeAttr(attr map[string]string) {
	for k, v := range attr {
		w.writeByte(' ')
		w.writeString(k)
		w.writeString(`="`)
		w.writeString(v)
		w.writeByte('"')
	}
}

// Bytes 将内容转换成 []byte 并返回
func (w *XMLWriter) Bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	return w.buf.Bytes(), nil
}

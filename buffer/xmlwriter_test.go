// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package buffer

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestWriter_writePI(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := &xmlWriter{
			buf: new(bytes.Buffer),
		}

		w.writePI(name, kv)
		bs, err := w.bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml-stylesheet", nil, "<?xml-stylesheet?>\n")
	test("xml-stylesheet", map[string]string{"type": "text/xsl"}, `<?xml-stylesheet type="text/xsl"?>`+"\n")
}

func TestWriter_writeElement(t *testing.T) {
	a := assert.New(t)
	test := func(name, val string, kv map[string]string, want string) {
		w := &xmlWriter{
			buf: new(bytes.Buffer),
		}

		w.writeElement(name, val, kv)
		bs, err := w.bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml", "text", nil, `<xml>text</xml>`+"\n")
	test("xml", "", nil, `<xml></xml>`+"\n")
	test("xml", "text", map[string]string{"type": "text/xsl"}, `<xml type="text/xsl">text</xml>`+"\n")
}

func TestWriter_writeCloseElement(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := &xmlWriter{
			buf: new(bytes.Buffer),
		}

		w.writeCloseElement(name, kv)
		bs, err := w.bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml", nil, `<xml />`+"\n")
	test("xml", nil, `<xml />`+"\n")
	test("xml", map[string]string{"type": "text/xsl"}, `<xml type="text/xsl" />`+"\n")
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package helper

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestWriter_WritePI(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := &XMLWriter{
			buf: new(bytes.Buffer),
		}

		w.WritePI(name, kv)
		bs, err := w.Bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml-stylesheet", nil, "<?xml-stylesheet?>\n")
	test("xml-stylesheet", map[string]string{"type": "text/xsl"}, `<?xml-stylesheet type="text/xsl"?>`+"\n")
}

func TestWriter_WriteElement(t *testing.T) {
	a := assert.New(t)
	test := func(name, val string, kv map[string]string, want string) {
		w := &XMLWriter{
			buf: new(bytes.Buffer),
		}

		w.WriteElement(name, val, kv)
		bs, err := w.Bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml", "text", nil, `<xml>text</xml>`+"\n")
	test("xml", "", nil, `<xml></xml>`+"\n")
	test("xml", "text", map[string]string{"type": "text/xsl"}, `<xml type="text/xsl">text</xml>`+"\n")
}

func TestWriter_WriteCloseElement(t *testing.T) {
	a := assert.New(t)
	test := func(name string, kv map[string]string, want string) {
		w := &XMLWriter{
			buf: new(bytes.Buffer),
		}

		w.WriteCloseElement(name, kv)
		bs, err := w.Bytes()
		a.NotError(err).Equal(string(bs), want)
	}

	test("xml", nil, `<xml />`+"\n")
	test("xml", nil, `<xml />`+"\n")
	test("xml", map[string]string{"type": "text/xsl"}, `<xml type="text/xsl" />`+"\n")
}

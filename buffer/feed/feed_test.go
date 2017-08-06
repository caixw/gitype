// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package feed

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

type pi struct {
	name string
	kv   map[string]string
	want string
}

func (pi *pi) test(a *assert.Assertion) {
	buf := new(bytes.Buffer)

	a.NotError(writePI(buf, pi.name, pi.kv))
	a.Equal(buf.String(), pi.want)
}

func TestWritePI(t *testing.T) {
	a := assert.New(t)

	(&pi{
		name: "xml-stylesheet",
		kv:   map[string]string{"type": "text/xsl"},
		want: `<?xml-stylesheet type="text/xsl" ?>`,
	}).test(a)

	(&pi{
		name: "xml-stylesheet",
		kv:   map[string]string{},
		want: `<?xml-stylesheet ?>`,
	}).test(a)
}

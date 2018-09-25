// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sw

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestServiceWorker_Add(t *testing.T) {
	a := assert.New(t)

	sw := New()
	a.NotNil(sw)

	sw.Add("v1", "p1", "p2")
	a.Equal(sw.caches["v1"], []string{"p1", "p2"})

	sw.Add("v1", "p3", "p4")
	a.Equal(sw.caches["v1"], []string{"p1", "p2", "p3", "p4"})

	sw.Add("v1", "p3", "p4")
	a.Equal(sw.caches["v1"], []string{"p1", "p2", "p3", "p4", "p3", "p4"})

	sw.Add("v2", "p1", "p2")
	a.Equal(sw.caches["v2"], []string{"p1", "p2"})
}

func TestServiceWorker_Bytes(t *testing.T) {
	a := assert.New(t)

	sw := New()
	a.NotNil(sw)

	sw.Add("v1", "p1", "p2")
	a.Equal(sw.caches["v1"], []string{"p1", "p2"})

	bs := sw.Bytes()
	a.True(bytes.Index(bs, replacement) < 0)
}

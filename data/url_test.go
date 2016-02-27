// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadURLS(t *testing.T) {
	a := assert.New(t)

	d := &Data{}
	a.NotError(d.loadURLS("./testdata/meta/urls.yaml"))
	urls := d.URLS
	a.NotNil(urls)
	a.Equal(urls.Post, "/posts").Equal(urls.Atom, "/atom.xml")
}

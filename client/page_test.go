// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPage_nextPage(t *testing.T) {
	a := assert.New(t)
	p := &page{}

	p.nextPage("url", "")
	a.Equal(p.NextPage.URL, "url")
	a.Equal(p.NextPage.Rel, "next")

	p.nextPage("url", "text")
	a.Equal(p.NextPage.URL, "url")
	a.Equal(p.NextPage.Rel, "next")
	a.Equal(p.NextPage.Text, "text")
}

func TestPage_prevPage(t *testing.T) {
	a := assert.New(t)
	p := &page{}

	p.prevPage("url", "")
	a.Equal(p.PrevPage.URL, "url")
	a.Equal(p.PrevPage.Rel, "prev")

	p.prevPage("url", "text")
	a.Equal(p.PrevPage.URL, "url")
	a.Equal(p.PrevPage.Rel, "prev")
	a.Equal(p.PrevPage.Text, "text")
}

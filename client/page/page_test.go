// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package page

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPage_Next(t *testing.T) {
	a := assert.New(t)
	p := &Page{}

	p.Next("url", "")
	a.Equal(p.NextPage.URL, "url")
	a.Equal(p.NextPage.Rel, "next")

	p.Next("url", "text")
	a.Equal(p.NextPage.URL, "url")
	a.Equal(p.NextPage.Rel, "next")
	a.Equal(p.NextPage.Text, "text")
}

func TestPage_Prev(t *testing.T) {
	a := assert.New(t)
	p := &Page{}

	p.Prev("url", "")
	a.Equal(p.PrevPage.URL, "url")
	a.Equal(p.PrevPage.Rel, "prev")

	p.Prev("url", "text")
	a.Equal(p.PrevPage.URL, "url")
	a.Equal(p.PrevPage.Rel, "prev")
	a.Equal(p.PrevPage.Text, "text")
}

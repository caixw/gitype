// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package helper

import (
	"testing"

	"github.com/issue9/assert"
)

func TestReplaceContent(t *testing.T) {
	a := assert.New(t)

	a.Equal("abc", ReplaceContent("abc%content%", ""))
	a.Equal("", ReplaceContent("", "TITLE"))
	a.Equal("TITLE", ReplaceContent("%content%", "TITLE"))
	a.Equal("abcTITLE", ReplaceContent("abc%content%", "TITLE"))
	a.Equal("abc%Content%", ReplaceContent("abc%Content%", "TITLE"))
}

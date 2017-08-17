// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/issue9/assert"
)

func TestAuthor_sanitize(t *testing.T) {
	a := assert.New(t)

	author := &Author{}
	a.Error(author.sanitize())

	author.Name = ""
	a.Error(author.sanitize())

	author.Name = "caixw"
	a.NotError(author.sanitize())
}

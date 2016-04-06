// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"sort"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &MetaError{}

func TestPostSort(t *testing.T) {
	a := assert.New(t)

	ps := []*Post{
		&Post{Slug: "4", Top: false, Created: 4},
		&Post{Slug: "2", Top: false, Created: 2},
		&Post{Slug: "3", Top: false, Created: 3},
		&Post{Slug: "1", Top: true, Created: 1},
		&Post{Slug: "0", Top: true, Created: 0},
	}

	sort.Sort(posts(ps))
	a.Equal(ps[0].Slug, "4")
	a.Equal(ps[1].Slug, "3")
	a.Equal(ps[2].Slug, "2")
	a.Equal(ps[3].Slug, "1")
	a.Equal(ps[4].Slug, "0")
}

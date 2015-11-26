// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestToMaps(t *testing.T) {
	a := assert.New(t)

	opt := &Options{
		PageSize: 30,
	}

	maps, err := opt.toMaps()
	a.NotError(err)
	for _, item := range maps {
		if item["group"] == "system" && item["key"] == "pageSize" {
			a.Equal(item["value"], "30")
		}
	}
}

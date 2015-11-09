// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package install

import (
	"testing"

	"github.com/caixw/typing/core"
	"github.com/issue9/assert"
)

func TestToMaps(t *testing.T) {
	a := assert.New(t)

	opt := &core.Options{
		PageSize: 30,
	}

	maps, err := toMaps(opt)
	a.NotError(err)
	for _, item := range maps {
		if item["group"] == "system" && item["key"] == "pageSize" {
			a.Equal(item["value"], "30")
		}
	}
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestLoad(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")
	d, err := Load(p)
	a.NotError(err).NotNil(d)
}

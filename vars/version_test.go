// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"testing"

	"github.com/issue9/assert"
	v "github.com/issue9/version"
)

func TestMainVersion(t *testing.T) {
	a := assert.New(t)

	a.True(v.SemVerValid(mainVersion))
}

// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sw

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestSWJS(t *testing.T) {
	a := assert.New(t)

	a.Equal(1, bytes.Count(swjs, replacement))
}

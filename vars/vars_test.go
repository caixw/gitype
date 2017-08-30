// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParseDate(t *testing.T) {
	a := assert.New(t)

	unix, err := ParseDate("2017-08-30T12:03:42+08:00")
	a.NotError(err).NotEmpty(unix)

	unix, err = ParseDate("2017-08-30")
	a.Error(err).Empty(unix)
}

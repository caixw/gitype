// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestHashPassword(t *testing.T) {
	a := assert.New(t)

	str1 := "123"
	str2 := HashPassword(str1)
	a.NotEmpty(str2).NotEqual(str1, str2)
}

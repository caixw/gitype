// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"fmt"
	"path"
	"testing"
)

func TestPathJoin(t *testing.T) {
	str := path.Join("abc/", "//def")
	fmt.Println(str)
}

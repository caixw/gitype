// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caixw/typing/app"
	"github.com/caixw/typing/path"
)

func main() {
	err := app.Run(path.New("./testdata"))

	if err != nil {
		panic(err)
	}
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/caixw/typing/admin"
	"github.com/caixw/typing/core"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/install"
	"github.com/caixw/typing/themes"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if install.Install() {
		return
	}

	err := core.Init()
	if err != nil {
		panic(err)
	}

	// themes
	if err = themes.Init(); err != nil {
		panic(err)
	}

	// admin
	if err := admin.Init(); err != nil {
		panic(err)
	}

	// feed
	if err = feed.Init(); err != nil {
		panic(err)
	}

	core.Run()
	core.Close()
}

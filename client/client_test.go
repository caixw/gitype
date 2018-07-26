// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"os"
	"testing"

	"github.com/issue9/web"
	"github.com/issue9/web/encoding"
	"github.com/issue9/web/encoding/html"

	"github.com/caixw/gitype/path"
)

var (
	client *Client
)

func TestMain(m *testing.M) {
	path := path.New("../testdata")
	var err error

	htmlMgr := html.New(nil)
	encoding.AddMarshal("text/html", htmlMgr.Marshal)

	client, err = New(path)
	if err != nil {
		panic(err)
	}

	if err = web.Init(path.ConfDir); err != nil {
		panic(err)
	}

	module := web.NewModule("test", "test")
	err = client.Mount(module.Mux(), htmlMgr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

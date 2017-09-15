// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
	"github.com/issue9/mux"
)

var (
	router = mux.New(false, false, nil, nil)
	server = httptest.NewServer(router)
	c      *Client
)

type httpTester struct {
	path    string
	content string
	status  int
}

func (t *httpTester) test(a *assert.Assertion) {
	// 非正常状态下，初始化 content 内容
	if len(t.content) == 0 && t.status > 299 {
		t.content = http.StatusText(t.status) + "\n"
	}

	resp, err := http.Get(server.URL + t.path)
	a.NotError(err).NotNil(resp)

	a.Equal(resp.StatusCode, t.status, "v1:%v,v2:%v,path:%v", resp.StatusCode, t.status, t.path)

	bs, err := ioutil.ReadAll(resp.Body)
	a.NotError(err).NotNil(bs)
	a.NotError(resp.Body.Close())

	if len(t.content) > 0 {
		a.Equal(bs, []byte(t.content), "v1:%v,v2:%v,path:%v", string(bs), t.content, t.path)
	}
}

func runHTTPTester(testers []*httpTester, t *testing.T) {
	a := assert.New(t)

	for _, test := range testers {
		test.test(a)
	}
}

func TestMain(t *testing.T) {
	a := assert.New(t)
	path := vars.NewPath("../testdata")

	client, err := New(path, router)
	a.NotError(err).NotNil(client)

	a.Equal(client.path, path)
	a.NotNil(client.data)
	a.Equal(client.Created(), client.data.Created)

	c = client
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

func TestRenderJSON(t *testing.T) {
	a := assert.New(t)

	w := httptest.NewRecorder()
	RenderJSON(w, http.StatusOK, nil, nil)
	a.Equal(w.Code, http.StatusOK).Equal(w.Body.String(), "")

	w = httptest.NewRecorder()
	RenderJSON(w, http.StatusInternalServerError, map[string]string{"name": "name"}, map[string]string{"h": "h"})
	a.Equal(w.Body.String(), `{"name":"name"}`)
	a.Equal(w.Header().Get("h"), "h")

	// 解析json出错，会返回500错误
	w = httptest.NewRecorder()
	RenderJSON(w, http.StatusOK, complex(5, 7), nil)
	a.Equal(w.Code, http.StatusInternalServerError)
	a.Equal(w.Body.String(), "")
}

func TestFileExists(t *testing.T) {
	a := assert.New(t)

	a.True(FileExists("util.go"))
	a.False(FileExists("unknown.go"))
}

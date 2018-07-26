// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"testing"

	"github.com/issue9/assert/rest"
	"github.com/issue9/web"
)

func TestGetAsset(t *testing.T) {
	h, err := web.Handler()
	if err != nil {
		panic(err)
	}
	s := rest.NewServer(t, h, nil)

	// getAsset
	s.NewRequest(http.MethodGet, "/posts/folder/post2/assets/assets.txt").
		Do().
		StringBody("assets.txt\n").
		Status(http.StatusOK)

	// content.html
	// 此条会优先匹配 getPost，然后跳转到 getRaw
	s.NewRequest(http.MethodGet, "/posts/folder/post2/content.html").
		Do().
		Status(http.StatusNotFound)

	// meta.yaml
	s.NewRequest(http.MethodGet, "/posts/folder/post2/meta.yaml").
		Do().
		Status(http.StatusNotFound)

	// 跳转到 getRaws
	s.NewRequest(http.MethodGet, "/posts/folder/post2/raws.txt").
		Do().
		StringBody("raws.txt\n").
		Status(http.StatusOK)
}

func TestGetTheme(t *testing.T) {
	h, err := web.Handler()
	if err != nil {
		panic(err)
	}
	s := rest.NewServer(t, h, nil)

	// css
	s.NewRequest(http.MethodGet, "/themes/t1/style.css").
		Do().
		StringBody("*{}\n").
		Status(http.StatusOK)

	// 模板文件
	s.NewRequest(http.MethodGet, "/themes/t1/template.html").
		Do().
		Status(http.StatusNotFound)

	// theme.yaml
	s.NewRequest(http.MethodGet, "/themes/t1/theme.yaml").
		Do().
		Status(http.StatusNotFound)

	// themes/analytics.html
	s.NewRequest(http.MethodGet, "/themes/analytics").
		Do().
		Status(http.StatusNotFound)

	// 跳转到 getRaws
	s.NewRequest(http.MethodGet, "/themes/t1/raws.txt").
		Do().
		StringBody("raws.txt\n").
		Status(http.StatusOK)
}

func TestGetRaws(t *testing.T) {
	h, err := web.Handler()
	if err != nil {
		panic(err)
	}
	s := rest.NewServer(t, h, nil)

	s.NewRequest(http.MethodGet, "/raws.txt").
		Do().
		StringBody("raws.txt\n").
		Status(http.StatusOK)

	s.NewRequest(http.MethodGet, "/not-exists.txt").
		Do().
		Status(http.StatusNotFound)
}

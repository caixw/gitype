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

func TestPost(t *testing.T) {
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

	// getPost
	s.NewRequest(http.MethodGet, "/posts/folder/post2.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// 跳转到 getRaws
	s.NewRequest(http.MethodGet, "/posts/folder/post2/raws.txt").
		Do().
		StringBody("raws.txt\n").
		Status(http.StatusOK)

	// getPosts，首页
	s.NewRequest(http.MethodGet, "/index.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// getPosts，第一页，肯定存在
	s.NewRequest(http.MethodGet, "/index.html?page=1").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// getPosts，肯定不存在
	s.NewRequest(http.MethodGet, "/index.html?page=100000").
		Do().
		Status(http.StatusNotFound)

	// getPosts，页码小于0
	s.NewRequest(http.MethodGet, "/index.html?page=-1").
		Do().
		Status(http.StatusNotFound)

	// getPosts
	s.NewRequest(http.MethodGet, "/").
		Do().
		BodyNotNil().
		Status(http.StatusOK)
}

func TestRoutes(t *testing.T) {
	h, err := web.Handler()
	if err != nil {
		panic(err)
	}
	s := rest.NewServer(t, h, nil)

	// archives.html
	s.NewRequest(http.MethodGet, "/archives.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// links.html
	s.NewRequest(http.MethodGet, "/links.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// tags.html
	s.NewRequest(http.MethodGet, "/tags.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// tags/default1.html?page=...
	s.NewRequest(http.MethodGet, "/tags/default1.html?page=1").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// tags/default1.html?page=...
	s.NewRequest(http.MethodGet, "/tags/default1.html?page=10000").
		Do().
		Status(http.StatusNotFound)

	// tags/default1.html?page=...
	s.NewRequest(http.MethodGet, "/tags/default1.html?page=0").
		Do().
		Status(http.StatusNotFound)

	// tags/...
	s.NewRequest(http.MethodGet, "/tags/default1.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// tags/...
	s.NewRequest(http.MethodGet, "/tags/default2.html").
		Do().
		BodyNotNil().
		Status(http.StatusOK)

	// tags/...
	s.NewRequest(http.MethodGet, "/tags/not-exists.html").
		Do().
		Status(http.StatusNotFound)

		// tags/...不存在并跳转到 getRaws
	s.NewRequest(http.MethodGet, "/tags/raws.html").
		Do().
		StringBody("raws.html\n").
		Status(http.StatusOK)
}

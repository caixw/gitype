// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"testing"
)

func TestPost(t *testing.T) {
	testers := []*httpTester{
		// getAsset
		{
			path:    "/posts/folder/post2/assets/assets.txt",
			content: "assets.txt\n",
			status:  http.StatusOK,
		},

		// getPost
		{
			path:   "/posts/folder/post2.html",
			status: http.StatusOK,
		},

		// 跳转到 getRaws
		{
			path:    "/posts/folder/post2/raws.txt",
			content: "raws.txt\n",
			status:  http.StatusOK,
		},

		// getPosts，首页
		{
			path:   "/index.html",
			status: http.StatusOK,
		},

		// getPosts，第一页，肯定存在
		{
			path:   "/index.html?page=1",
			status: http.StatusOK,
		},

		// getPosts，不存在的页码
		{
			path:   "/index.html?page=10000",
			status: http.StatusNotFound,
		},

		// getPosts，页码小于0
		{
			path:   "/index.html?page=-1",
			status: http.StatusNotFound,
		},

		// getPosts
		{
			path:   "/",
			status: http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

func TestRoutes(t *testing.T) {
	testers := []*httpTester{
		// archives.html
		{
			path:   "/archives.html",
			status: http.StatusOK,
		},
		// links.html
		{
			path:   "/links.html",
			status: http.StatusOK,
		},
		// tags.html
		{
			path:   "/tags.html",
			status: http.StatusOK,
		},

		// tags.html?page=...
		{
			path:   "/tags/default1.html?page=1",
			status: http.StatusOK,
		},

		// tags.html?page=...
		{
			path:   "/tags/default1.html?page=10000",
			status: http.StatusNotFound,
		},

		// tags.html?page=...
		{
			path:   "/tags/default1.html?page=0",
			status: http.StatusNotFound,
		},

		// tags/...
		{
			path:   "/tags/default1.html",
			status: http.StatusOK,
		},

		// tags/...
		{
			path:   "/tags/default2.html",
			status: http.StatusOK,
		},

		// tags/...
		{
			path:   "/tags/not-exists.html",
			status: http.StatusNotFound,
		},

		// tags/...不存在并跳转到 getRaws
		{
			path:    "/tags/raws.html",
			content: "raws.html\n",
			status:  http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

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
		&httpTester{
			path:    "/posts/folder/post2/assets/assets.txt",
			content: "assets.txt\n",
			status:  http.StatusOK,
		},

		// getPost
		&httpTester{
			path:   "/posts/folder/post2.html",
			status: http.StatusOK,
		},

		// 跳转到 getRaws
		&httpTester{
			path:    "/posts/folder/post2/raws.txt",
			content: "raws.txt\n",
			status:  http.StatusOK,
		},

		// getPosts，首页
		&httpTester{
			path:   "/index.html",
			status: http.StatusOK,
		},

		// getPosts，第一页，肯定存在
		&httpTester{
			path:   "/index.html?page=1",
			status: http.StatusOK,
		},

		// getPosts，不存在的页码
		&httpTester{
			path:   "/index.html?page=10000",
			status: http.StatusNotFound,
		},

		// getPosts，页码小于0
		&httpTester{
			path:   "/index.html?page=-1",
			status: http.StatusNotFound,
		},

		// getPosts
		&httpTester{
			path:   "/",
			status: http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

func TestRoutes(t *testing.T) {
	testers := []*httpTester{
		// archives.html
		&httpTester{
			path:   "/archives.html",
			status: http.StatusOK,
		},
		// links.html
		&httpTester{
			path:   "/links.html",
			status: http.StatusOK,
		},
		// tags.html
		&httpTester{
			path:   "/tags.html",
			status: http.StatusOK,
		},

		// tags.html?page=...
		&httpTester{
			path:   "/tags/default1.html?page=1",
			status: http.StatusOK,
		},

		// tags.html?page=...
		&httpTester{
			path:   "/tags/default1.html?page=10000",
			status: http.StatusNotFound,
		},

		// tags.html?page=...
		&httpTester{
			path:   "/tags/default1.html?page=0",
			status: http.StatusNotFound,
		},

		// tags/...
		&httpTester{
			path:   "/tags/default1.html",
			status: http.StatusOK,
		},

		// tags/...
		&httpTester{
			path:   "/tags/default2.html",
			status: http.StatusOK,
		},

		// tags/...
		&httpTester{
			path:   "/tags/not-exists.html",
			status: http.StatusNotFound,
		},

		// tags/...不存在并跳转到 getRaws
		&httpTester{
			path:    "/tags/raws.html",
			content: "raws.html\n",
			status:  http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

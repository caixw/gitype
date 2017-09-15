// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestIsIgnoreThemeFile(t *testing.T) {
	a := assert.New(t)
	a.True(isIgnoreThemeFile(vars.TemplateExtension))
	a.True(isIgnoreThemeFile(".yaml"))
	a.False(isIgnoreThemeFile(".txt"))
	a.False(isIgnoreThemeFile(".css"))
}

func TestGetAsset(t *testing.T) {
	testers := []*httpTester{
		&httpTester{
			path:    "/posts/folder/post2/assets/assets.txt",
			content: "assets.txt\n",
			status:  http.StatusOK,
		},

		// content.html
		// 此条会优先匹配 getPost，然后跳转到 getRaw
		&httpTester{
			path:   "/posts/folder/post2/content.html",
			status: http.StatusNotFound,
		},

		// meta.yaml
		&httpTester{
			path:   "/posts/folder/post2/meta.yaml",
			status: http.StatusNotFound,
		},

		// 跳转到 getRaws
		&httpTester{
			path:    "/posts/folder/post2/raws.txt",
			content: "raws.txt\n",
			status:  http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

func TestGetTheme(t *testing.T) {
	testers := []*httpTester{
		&httpTester{
			path:    "/themes/t1/style.css",
			content: "*{}\n",
			status:  http.StatusOK,
		},

		// 模板文件
		&httpTester{
			path:   "/themes/t1/template.html",
			status: http.StatusNotFound,
		},

		// theme.yaml
		&httpTester{
			path:   "/themes/t1/theme.yaml",
			status: http.StatusNotFound,
		},

		// 跳转到 getRaws
		&httpTester{
			path:    "/themes/t1/raws.txt",
			content: "raws.txt\n",
			status:  http.StatusOK,
		},
	}

	runHTTPTester(testers, t)
}

func TestGetRaws(t *testing.T) {
	testers := []*httpTester{
		&httpTester{
			path:    "/raws.txt",
			content: "raws.txt\n",
			status:  http.StatusOK,
		},

		&httpTester{
			path:   "/not-exists.txt",
			status: http.StatusNotFound,
		},
	}

	runHTTPTester(testers, t)
}

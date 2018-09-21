// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadTheme(t *testing.T) {
	a := assert.New(t)

	theme, err := LoadTheme(testdataPath, "t1")
	a.NotError(err).NotNil(theme)

	a.Equal(theme.Name, "name")
	a.Equal(theme.Author.Name, "caixw")
}

func TestLoadLinks(t *testing.T) {
	a := assert.New(t)

	links, err := LoadLinks(testdataPath)
	a.NotError(err).NotNil(links)

	a.True(len(links) > 0)
	a.Equal(links[0].Text, "text0")
	a.Equal(links[0].URL, "url0")
	a.Equal(links[1].Text, "text1")
	a.Equal(links[1].URL, "url1")
}

func TestAuthor_sanitize(t *testing.T) {
	a := assert.New(t)

	author := &Author{}
	a.Error(author.sanitize())

	author.Name = ""
	a.Error(author.sanitize())

	author.Name = "caixw"
	a.NotError(author.sanitize())
}

func TestIcon_sanitize(t *testing.T) {
	a := assert.New(t)

	icon := &Icon{}
	a.Error(icon.sanitize())

	// type 会根据 URL 的扩展名，可能会自动计算获得。
	// png 有内置，肯定可以成功检测到
	icon = &Icon{URL: "http://example.com/1.png"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "image/png")

	// 不存在的扩展名，则不会计算工其 type
	icon = &Icon{URL: "http://example.com/1.not-exists"}
	a.NotError(icon.sanitize())
	a.Equal(icon.Type, "")
}

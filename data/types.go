// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"github.com/caixw/gitype/data/loader"
)

// Feed RSS、Atom、Sitemap 和 Opensearch 的配置内容
type Feed struct {
	Title   string // 标题，一般出现在 html>head>link.title 属性中
	URL     string // 地址，不能包含域名
	Type    string // mime type
	Content []byte // 实际的内容
}

// Author 描述作者信息
type Author = loader.Author

// Link 描述链接的内容
type Link = loader.Link

// Icon 表示网站图标，比如 html>head>link.rel="short icon"
type Icon = loader.Icon

// Page 配置配置
type Page = loader.Page

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client ...
package client

import (
	"html/template"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// Client 所有数据的缓存，每次更新数据时，
// 直接声明一个新的 Client 实例，丢弃原来的 Client 即可。
type Client struct {
	Created int64 // 当前数据的加载时间
	Data    *data.Data

	// 缓存的数据
	Template   *template.Template // 主题编译后的模板
	RSS        []byte
	Atom       []byte
	Sitemap    []byte
	Opensearch []byte
}

// New 声明一个新的 Buffer 实例
func New(path *vars.Path) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	b := &Client{
		Created: time.Now().Unix(),
		Data:    d,
	}

	errFilter := func(fn func() error) {
		if err == nil {
			err = fn()
		}
	}

	errFilter(b.compileTemplate)
	errFilter(b.buildRSS)
	errFilter(b.buildAtom)
	errFilter(b.buildSitemap)
	errFilter(b.buildOpensearch)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func formatUnix(unix int64, format string) string {
	t := time.Unix(unix, 0)
	return t.Format(format)
}

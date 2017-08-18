// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package buffer 加载数据以及对数据延伸内容的一些处理，比如根据配置文件生成 RSS 等内容。
//
// buffer 是对 data 包的增强，提供的数据依然是相对比较固定的，动态内容应该由 app 包来生成。
package buffer

import (
	"html/template"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// Buffer 所有数据的缓存，每次更新数据时，
// 直接声明一个新的 Buffer 实例，丢弃原来的 Buffer 即可。
type Buffer struct {
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
func New(path *vars.Path) (*Buffer, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	b := &Buffer{
		Created: time.Now().Unix(),
		Data:    d,
	}

	if err = b.compileTemplate(); err != nil {
		return nil, err
	}

	if err = b.initFeeds(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Buffer) initFeeds() error {
	conf := b.Data.Config

	if conf.RSS != nil {
		rss, err := buildRSS(b.Data)
		if err != nil {
			return err
		}
		b.RSS = rss
	}

	if conf.Atom != nil {
		atom, err := buildAtom(b.Data)
		if err != nil {
			return err
		}

		b.Atom = atom
	}

	if conf.Sitemap != nil {
		sitemap, err := buildSitemap(b.Data)
		if err != nil {
			return err
		}

		b.Sitemap = sitemap
	}

	if conf.Opensearch != nil {
		opensearch, err := buildOpensearch(b.Data)
		if err != nil {
			return err
		}

		b.Opensearch = opensearch
	}

	return nil
}

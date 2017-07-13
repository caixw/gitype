// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package buffer

import (
	"html/template"
	"strconv"
	"time"

	"github.com/caixw/typing/buffer/feed"
	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

// Buffer 所有数据的缓存
type Buffer struct {
	path *vars.Path

	Updated  int64              // 更新时间，一般为重新加载数据的时间
	Etag     string             // 所有页面都采用相同的 Etag，即时间戳字符串
	Data     *data.Data         // 加载的数据，每次加载都会被重置
	Template *template.Template // 主题编译后的模板

	RSS     []byte
	Atom    []byte
	Sitemap []byte
}

// New 声明一个新的 Buffer 实例
func New(path *vars.Path) (*Buffer, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	b := &Buffer{
		path:    path,
		Data:    d,
		Updated: now,
		Etag:    strconv.FormatInt(now, 10),
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
		rss, err := feed.BuildRSS(b.Data)
		if err != nil {
			return err
		}
		b.RSS = rss.Bytes()
	}

	if conf.Atom != nil {
		atom, err := feed.BuildAtom(b.Data)
		if err != nil {
			return err
		}

		b.Atom = atom.Bytes()
	}

	if conf.Sitemap != nil {
		sitemap, err := feed.BuildSitemap(b.Data)
		if err != nil {
			return err
		}

		b.Sitemap = sitemap.Bytes()
	}

	return nil
}

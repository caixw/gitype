// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 加载数据并对其进行处理。
package data

import (
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
	"golang.org/x/text/search"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path    *path.Path
	Created time.Time

	// Updated 数据的更新时间，诸如 outdatedServer 等服务，
	// 会定时更新数据，Updated 即记录这些更新的时间。
	Updated time.Time

	// Etag 表示 根据 Updated 生成的 etag 字符串
	Etag string

	// 直接从 config 中继承过来的变量
	SiteName string
	Subtitle string           // 网站副标题
	Language string           // 语言标记，比如 zh-cmn-Hans
	Beian    string           // 备案号
	Uptime   time.Time        // 上线时间
	PageSize int              // 每页显示的数量
	Type     string           // 页面的 mime type 类型
	Icon     *Icon            // 程序默认的图标
	Menus    []*Link          // 导航菜单
	Author   *Author          // 默认作者信息
	License  *Link            // 默认版权信息
	Pages    map[string]*Page // 各个页面的自定义内容

	outdatedServer *outdatedServer

	Tags     []*Tag
	Series   []*Tag
	Links    []*Link
	Posts    []*Post
	Archives []*Archive
	Theme    *Theme // 当前主题

	Opensearch *Feed
	Sitemap    *Feed
	RSS        *Feed
	Atom       *Feed

	Matcher *search.Matcher
}

// Load 函数用于加载一份新的数据。
func Load(path *path.Path) (*Data, error) {
	conf, err := loader.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	tags, err := loadTags(path, conf)
	if err != nil {
		return nil, err
	}

	links, err := loader.LoadLinks(path)
	if err != nil {
		return nil, err
	}

	posts, err := loadPosts(path, tags, conf)
	if err != nil {
		return nil, err
	}

	theme, err := loadTheme(path, conf)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	d := &Data{
		path:    path,
		Created: now,

		SiteName: conf.Title,
		Language: conf.Language,
		Subtitle: conf.Subtitle,
		Beian:    conf.Beian,
		Uptime:   conf.Uptime,
		PageSize: conf.PageSize,
		Type:     conf.Type,
		Icon:     conf.Icon,
		Menus:    conf.Menus,
		Pages:    conf.Pages,

		Tags:  tags,
		Links: links,
		Posts: posts,
		Theme: theme,

		Matcher: search.New(conf.LanguageTag, search.Loose),
	}

	if err := d.sanitize(conf); err != nil {
		return nil, err
	}

	d.initOutdatedServer(conf)

	d.setUpdated(now)
	return d, nil
}

// Free 释放数据内容
func (d *Data) Free() {
	d.outdatedServer.stop()
}

// 调整更新时间
func (d *Data) setUpdated(t time.Time) {
	d.Updated = t
	d.Etag = vars.Etag(t)
}

// 对各个数据再次进行检测，主要是一些关联数据的相互初始化
func (d *Data) sanitize(conf *loader.Config) error {
	if err := d.compileTemplate(); err != nil {
		return err
	}

	// 过滤空标签
	tags := make([]*Tag, 0, len(d.Tags))
	for _, tag := range d.Tags {
		if len(tag.Posts) == 0 {
			continue
		}

		tags = append(tags, tag)
	}

	// 最后才分离标签和专题
	d.Tags, d.Series = splitTags(tags)

	return d.buildData(conf)
}

func (d *Data) buildData(conf *loader.Config) (err error) {
	errFilter := func(fn func(*loader.Config) error) {
		if err != nil {
			return
		}
		err = fn(conf)
	}

	errFilter(d.buildArchives)
	errFilter(d.buildOpensearch)
	errFilter(d.buildSitemap)
	errFilter(d.buildRSS)
	errFilter(d.buildAtom)
	return err
}

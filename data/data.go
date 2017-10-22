// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 加载数据并对其进行处理。
package data

import (
	"strings"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path *path.Path

	// Created 当前实例的加载时间，当数据加载完成之后，
	// 此值不再变化。
	Created time.Time

	// Updated 数据实例中的数据最后修改时间，
	// Data 内部可能会修改数据，比如每天会更新文章的过期提醒内容。
	Updated time.Time

	// Etag 用于网页请求时返回的 Etag 报头，根据 Updated 字段获取。
	Etag string

	// 直接从 config 中继承过来的变量
	SiteName string
	Subtitle string           // 网站副标题
	Language string           // 语言标记，比如 zh-cmn-Hans
	URL      string           // 网站的域名
	Beian    string           // 备案号
	Uptime   time.Time        // 上线时间
	PageSize int              // 每页显示的数量
	Type     string           // 页面的 mime type 类型
	Icon     *Icon            // 程序默认的图标
	Menus    []*Link          // 导航菜单
	Author   *Author          // 默认作者信息
	License  *Link            // 默认版权信息
	Pages    map[string]*Page // 各个页面的自定义内容

	outdated        *outdatedConfig
	postsTicker     *time.Ticker // 用于更新文章 Outdated 属性的定时器
	postsTickerDone chan bool

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
}

// Load 函数用于加载一份新的数据。
func Load(path *path.Path) (*Data, error) {
	conf, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	tags, err := loadTags(path)
	if err != nil {
		return nil, err
	}

	links, err := loadLinks(path)
	if err != nil {
		return nil, err
	}

	posts, err := loadPosts(path)
	if err != nil {
		return nil, err
	}

	theme, err := findTheme(path, conf)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	d := &Data{
		path:    path,
		Created: now,
		Updated: now,
		Etag:    vars.Etag(now),

		SiteName: conf.Title,
		Language: conf.Language,
		Subtitle: conf.Subtitle,
		URL:      conf.URL,
		Beian:    conf.Beian,
		Uptime:   conf.Uptime,
		PageSize: conf.PageSize,
		Type:     conf.Type,
		Icon:     conf.Icon,
		Menus:    conf.Menus,
		Pages:    conf.Pages,

		outdated:        conf.Outdated,
		postsTicker:     time.NewTicker(day),
		postsTickerDone: make(chan bool, 1),

		Tags:  tags,
		Links: links,
		Posts: posts,
		Theme: theme,
	}

	if err := d.sanitize(conf); err != nil {
		return nil, err
	}

	if err := d.buildData(conf); err != nil {
		return nil, err
	}

	// 所有数据没问题，则运行 ticker
	d.runUpdateOutdatedServer()

	return d, nil
}

// Free 释放当前数据内容
func (d *Data) Free() {
	d.stopPostsTicker()
}

// 对各个数据再次进行检测，主要是一些关联数据的相互初始化
func (d *Data) sanitize(conf *config) error {
	if err := d.compileTemplate(); err != nil {
		return err
	}

	p := conf.Pages[vars.PageTag]
	for _, tag := range d.Tags {
		// 将标签的默认修改时间设置为网站的上线时间
		tag.Modified = conf.Uptime

		tag.HTMLTitle = helper.ReplaceContent(p.Title, tag.Title)
	}

	for _, post := range d.Posts {
		if post.Author == nil {
			post.Author = conf.Author
		}

		if post.License == nil {
			post.License = conf.License
		}

		if err := d.attachPostTag(post, conf); err != nil {
			return err
		}
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
	ts, series := splitTags(tags)
	d.Tags = ts
	d.Series = series

	return nil
}

// 关联文章与标签的相关信息
func (d *Data) attachPostTag(post *Post, conf *config) *helper.FieldError {
	ts := strings.Split(post.TagsString, ",")
	for _, tag := range d.Tags {
		for _, slug := range ts {
			if tag.Slug != slug {
				continue
			}

			post.Tags = append(post.Tags, tag)
			tag.Posts = append(tag.Posts, post)

			if tag.Modified.Before(post.Modified) {
				tag.Modified = post.Modified
			}
			break
		}
	} // end for tags

	post.HTMLTitle = helper.ReplaceContent(conf.Pages[vars.PagePost].Title, post.Title)

	if len(post.Tags) == 0 {
		return &helper.FieldError{File: d.path.PostMetaPath(post.Slug), Message: "未指定任何关联标签信息", Field: "tags"}
	}

	return nil
}

func (d *Data) buildData(conf *config) (err error) {
	errFilter := func(fn func(*config) error) {
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

// BuildURL 生成一个带域名的地址
func (d *Data) BuildURL(path string) string {
	return d.URL + path
}

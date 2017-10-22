// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 加载数据并对其进行处理。
package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path    *path.Path
	Created time.Time

	Title    string           // 网站标题
	Language string           // 语言标记，比如 zh-cmn-Hans
	Subtitle string           // 网站副标题
	URL      string           // 网站的域名，非默认端口也得包含，不包含最后的斜杠，仅在生成地址时使用
	Beian    string           // 备案号
	Uptime   time.Time        // 上线时间
	PageSize int              // 每页显示的数量
	Type     string           // 所有页面的 mime type 类型，默认使用
	Icon     *Icon            // 程序默认的图标
	Menus    []*Link          // 导航菜单
	Author   *Author          // 默认作者信息
	License  *Link            // 默认版权信息
	Pages    map[string]*Page // 各个页面的自定义内容

	longDateFormat  string // 长时间的显示格式
	shortDateFormat string // 短时间的显示格式
	outdated        *outdatedConfig

	Tags     []*Tag
	Series   []*Tag
	Links    []*Link
	Posts    []*Post
	Archives []*Archive
	Themes   []*Theme // 所有可用主题，第一个元素为默认主题

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

	themes, err := loadThemes(path)
	if err != nil {
		return nil, err
	}

	d := &Data{
		path:    path,
		Created: time.Now(),

		Title:    conf.Title,
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

		longDateFormat:  conf.LongDateFormat,
		shortDateFormat: conf.ShortDateFormat,
		outdated:        conf.Outdated,

		Tags:   tags,
		Links:  links,
		Posts:  posts,
		Themes: themes,
	}

	if err := d.sanitize(conf); err != nil {
		return nil, err
	}

	if err := d.buildData(conf); err != nil {
		return nil, err
	}

	return d, nil
}

// 对各个数据再次进行检测，主要是一些关联数据的相互初始化
func (d *Data) sanitize(conf *config) error {
	if err := d.sanitizeThemes(conf); err != nil {
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

// Outdated 计算指定文章的 Outdated 信息。
// Outdated 是一个动态的值（其中的天数会变化），必须是在请求时生成。
func (d *Data) Outdated(post *Post) {
	if d.outdated == nil {
		return
	}

	now := time.Now()
	var outdated time.Duration

	switch d.outdated.Type {
	case outdatedTypeCreated:
		outdated = now.Sub(post.Created)
	case outdatedTypeModified:
		outdated = now.Sub(post.Modified)
	default:
		// 理论上此段代码永远不会运行，除非代码中直接修改了 Data.outdated.type 的值，
		// 因为在 outdatedConfig.sanitize 中已经作了判断。
		panic("无效的 config.yaml/outdated.type")
	}

	if outdated >= d.outdated.Duration {
		post.Outdated = fmt.Sprintf(d.outdated.Content, int64(outdated.Hours())/24)
	}
}

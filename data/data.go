// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 加载和再加工所有数据
package data

import (
	"html/template"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/caixw/typing/vars"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path    *vars.Path
	Created time.Time

	Config   *Config
	Theme    *Theme             // 当前主题
	Template *template.Template // 当前主题的模板
	Themes   []*Theme           // 主题列表
	Tags     []*Tag
	Links    []*Link
	Posts    []*Post
	Archives []*Archive

	Opensearch *Opensearch
	Sitemap    *Sitemap
	RSS        *RSS
	Atom       *RSS
}

// Load 函数用于加载一份新的数据。
func Load(path *vars.Path) (*Data, error) {
	// conf 需要先初始化
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
		Config:  newConfig(conf),
		Tags:    tags,
		Links:   links,
		Posts:   posts,
		Themes:  themes,
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
	for _, theme := range d.Themes { // 检测配置文件中的主题是否存在
		if theme.ID == conf.Theme {
			d.Theme = theme
			break
		}
	}
	if d.Theme == nil {
		return &FieldError{File: d.path.MetaConfigFile, Message: "该主题并不存在", Field: "theme"}
	}

	for _, tag := range d.Tags { // 将标签的默认修改时间设置为网站的上线时间
		tag.Modified = conf.Uptime
	}

	if err := d.attachPostMeta(conf); err != nil {
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
	d.Tags = tags

	return nil
}

// 关联文章的相关属性
func (d *Data) attachPostMeta(conf *config) *FieldError {
	for _, post := range d.Posts {
		if post.Author == nil {
			post.Author = conf.Author
		}

		if post.License == nil {
			post.License = conf.License
		}

		// tags
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

		if len(post.Tags) == 0 {
			return &FieldError{File: post.Slug, Message: "未指定任何关联标签信息", Field: "tags"}
		}
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
	if err != nil {
		return err
	}

	return d.compileTemplate()
}

// 加载 yaml 格式的文件 path 中的内容到 obj
func loadYamlFile(path string, obj interface{}) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, obj)
}

func (d *Data) url(path string) string {
	return d.Config.URL + path
}

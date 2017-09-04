// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 加载和再加工所有数据
package data

import (
	"html/template"
	"io/ioutil"
	"strconv"
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
	conf := &config{}
	if err := loadYamlFile(path.MetaConfigFile, conf); err != nil {
		return nil, err
	}
	if err := conf.sanitize(); err != nil {
		err.Field = path.MetaConfigFile
		return nil, err
	}

	d := &Data{
		path:    path,
		Created: time.Now(),
		Config:  newConfig(conf),
	}

	if err := d.loadFiles(); err != nil {
		return nil, err
	}

	if err := d.sanitize(conf); err != nil {
		return nil, err
	}

	if err := d.sanitize2(conf); err != nil {
		return nil, err
	}

	if err := d.buildData(conf); err != nil {
		return nil, err
	}

	return d, nil
}

// 加载所有的文件
func (d *Data) loadFiles() error {
	tags := make([]*Tag, 0, 100)
	if err := loadYamlFile(d.path.MetaTagsFile, &tags); err != nil {
		return err
	}
	d.Tags = tags

	links := make([]*Link, 0, 20)
	if err := loadYamlFile(d.path.MetaLinksFile, &links); err != nil {
		return err
	}
	d.Links = links

	posts, err := loadPosts(d.path)
	if err != nil {
		return err
	}
	d.Posts = posts

	themes, err := loadThemes(d.path)
	if err != nil {
		return err
	}
	d.Themes = themes

	return nil
}

// 对各个加载的数据进行转换、审查等操作。
func (d *Data) sanitize(conf *config) error {
	for index, tag := range d.Tags {
		if err := tag.sanitize(); err != nil {
			err.File = d.path.MetaTagsFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	for index, link := range d.Links {
		if err := link.sanitize(); err != nil {
			err.File = d.path.MetaLinksFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	for _, theme := range d.Themes {
		if err := theme.sanitize(); err != nil {
			err.File = theme.Path
			return err
		}
	}

	for _, post := range d.Posts {
		if err := post.sanitize(); err != nil {
			return err
		}
	}

	return nil
}

// 对各个数据再次进行检测，主要是一些关联数据的相互初始化
func (d *Data) sanitize2(conf *config) error {
	// 对文章进行排序，需保证 created 已经被初始化
	sortPosts(d.Posts)

	// 检测配置文件中的主题是否存在
	for _, theme := range d.Themes {
		if theme.ID == conf.Theme {
			d.Theme = theme
			break
		}
	}
	if d.Theme == nil {
		return &FieldError{File: d.path.MetaConfigFile, Message: "该主题并不存在", Field: "theme"}
	}

	// 将标签的修改时间设置为网站的上线时间
	for _, tag := range d.Tags {
		tag.Modified = conf.Uptime
	}

	if err := d.attachPostMeta(conf); err != nil {
		return err
	}

	// 过滤空标签，排序标签关联的文章
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
		if len(ts) == 0 {
			return &FieldError{File: post.Slug, Message: "未指定任何关联标签信息", Field: "tags"}
		}
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

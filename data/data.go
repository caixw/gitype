// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 负责加载 data 目录下的数据，以及一些固有格式的转换，比如时间格式。
package data

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/caixw/typing/vars"
)

const (
	tagsFilename  = "tags.yaml"
	linksFilename = "links.yaml"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path   *vars.Path
	Config *Config  // 配置内容
	Theme  *Theme   // 当前主题
	Tags   []*Tag   // map 对顺序是未定的，所以使用 slice
	Links  []*Link  // 友情链接
	Posts  []*Post  // 所有的文章列表
	Themes []*Theme // 主题，使用 slice，方便排序
}

// Load 函数用于加载一份新的数据。
func Load(path *vars.Path) (*Data, error) {
	d := &Data{
		path: path,
	}

	if err := d.loadFiles(); err != nil {
		return nil, err
	}

	if err := d.sanitize(); err != nil {
		return nil, err
	}

	if err := d.sanitize2(); err != nil {
		return nil, err
	}

	return d, nil
}

// 加载所有的文件
func (d *Data) loadFiles() error {
	tags := make([]*Tag, 0, 100)
	if err := loadYamlFile(d.path.MetaPath(tagsFilename), &tags); err != nil {
		return err
	}
	d.Tags = tags

	links := make([]*Link, 0, 20)
	if err := loadYamlFile(d.path.MetaPath(linksFilename), &links); err != nil {
		return err
	}
	d.Links = links

	config := &Config{}
	if err := loadYamlFile(d.path.MetaPath(confFilename), config); err != nil {
		return err
	}
	d.Config = config

	posts, err := loadPosts(d.path.PostsDir)
	if err != nil {
		return err
	}
	d.Posts = posts

	themes, err := loadThemes(d.path.ThemesDir)
	if err != nil {
		return err
	}
	d.Themes = themes

	return nil
}

// 对各个加载的数据进行转换、审查等操作。
func (d *Data) sanitize() error {
	for index, tag := range d.Tags {
		if err := tag.sanitize(); err != nil {
			err.File = tagsFilename
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	for index, link := range d.Links {
		if err := link.sanitize(); err != nil {
			err.File = linksFilename
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	if err := d.Config.sanitize(); err != nil {
		return err
	}

	for _, theme := range d.Themes {
		if err := theme.sanitize(); err != nil {
			err.File = filepath.Join(theme.Path, theme.ID)
			return err
		}
	}

	for _, post := range d.Posts {
		if err := post.sanitize(); err != nil {
			return err
		}
	}
	sortPosts(d.Posts)

	return nil
}

// 对各个数据再次进行检测，主要是一些关联数据的相互初始化
func (d *Data) sanitize2() error {
	// 检测配置文件中的主题是否存在
	for _, theme := range d.Themes {
		if theme.ID == d.Config.Theme {
			d.Theme = theme
			break
		}
	}
	if d.Theme == nil {
		return &FieldError{File: confFilename, Message: "该主题并不存在", Field: "theme"}
	}

	// 将标签的修改时间设置为网站的上线时间
	for _, tag := range d.Tags {
		tag.Modified = d.Config.Uptime
	}

	// 关联文章与标签
	if err := d.attachPostTags(); err != nil {
		return err
	}

	// 过滤空标签，排序标签关联的文章
	tags := make([]*Tag, 0, len(d.Tags))
	for _, tag := range d.Tags {
		if len(tag.Posts) == 0 {
			continue
		}

		sortPosts(tag.Posts)
		tags = append(tags, tag)
	}
	d.Tags = tags

	return nil
}

func (d *Data) attachPostTags() *FieldError {
	for _, post := range d.Posts {
		if post.Author == nil {
			post.Author = d.Config.Author
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

				if tag.Modified < post.Modified {
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

// 加载 yaml 格式的文件 path 中的内容到 obj
func loadYamlFile(path string, obj interface{}) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bs, obj)
}

// 分析时间，将其转换成 unix 时间戳
func parseDate(format string) (int64, error) {
	t, err := time.Parse(vars.DateFormat, format)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

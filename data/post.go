// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/issue9/logs"

	"gopkg.in/yaml.v2"
)

type Post struct {
	Slug           string  `yaml:"slug"`     // 唯一名称
	Title          string  `yaml:"title"`    // 标题
	Created        int64   `yaml:"-"`        // 创建时间
	Modified       int64   `yaml:"-"`        // 修改时间
	Tags           []*Tag  `yaml:"-"`        // 关联的标签
	Author         *Author `yaml:"author"`   // 作者
	Template       string  `yaml:"template"` // 使用的模板。未指定，则使用系统默认的
	Top            bool    `yaml:"top"`      // 是否置顶，多个置顶，则按时间排序
	Summary        string  `yaml:"summary"`  // 摘要
	Content        string  `yaml:"-"`        // 内容
	CreatedFormat  string  `yaml:"created"`  // 创建时间的字符串表示形式
	ModifiedFormat string  `yaml:"modified"` // 修改时间的字符串表示形式
	TagsString     string  `yaml:"tags"`     // 关联标签的列表
	Path           string  `yaml:"path"`     // 正文的文件名，相对于meta所在的目录
}

// 查找指定名称的文章。
// 若返回nil，则表示该文章不存在。
func (d *Data) FindPost(slug string) *Post {
	for _, post := range d.Posts {
		if post.Slug == slug {
			return post
		}
	}

	return nil
}

// 加载所有的文章内容。
// dir data/posts目录。
func (d *Data) loadPosts() error {
	dir := filepath.Join(d.path, "posts")

	paths := make([]string, 0, 100)

	// 遍历data/posts目录，查找所有的meta.yaml文章。
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "meta.yaml" {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return err
	}

	// 开始加载文章的具体内容。
	d.Posts = make([]*Post, 0, len(paths))
	for _, path := range paths {
		p, err := loadPost(dir, path, d.Config, d.Tags)
		if err != nil {
			logs.Error(err)
			continue
		}

		d.Posts = append(d.Posts, p)
	}

	return nil
}

// 加载某一文章的元数据。不包含实际内容。
// postsDir 表示data/posts目录的绝对地址；
// path 表示具体文章的meta.yaml文章；
func loadPost(postsDir, path string, conf *Config, tags []*Tag) (*Post, error) {
	dir := filepath.Dir(path)                 // 获取路径部分
	name := strings.TrimPrefix(dir, postsDir) // 获取相对于data/posts 的名称

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := &Post{}
	if err := yaml.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("[%v]解板yaml出错:%v", name, err)
	}

	if len(p.Slug) == 0 {
		return nil, fmt.Errorf("[%v]:文章唯一名称不能为空", name)

	}

	if len(p.Title) == 0 {
		return nil, fmt.Errorf("[%v]:文章标题不能为空", name)
	}

	if p.Author == nil {
		p.Author = conf.Author
	}

	// content
	if len(p.Path) == 0 {
		return nil, fmt.Errorf("[%v]:未指定内容文件", name)
	}
	data, err = ioutil.ReadFile(filepath.Join(dir, p.Path))
	if err != nil {
		return nil, fmt.Errorf("[%v]:读取文章内容出错：[%v]", name, err)
	}
	p.Content = string(data)

	// tags
	ts := strings.Split(p.TagsString, ",")
	for _, tag := range tags {
		for _, tagName := range ts {
			if tag.Slug == tagName {
				p.Tags = append(p.Tags, tag)
				break
			}
		}
	}

	// created
	t, err := time.Parse(parseDateFormat, p.CreatedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其创建时间是出错：[%v]", name, err)
	}
	p.Created = t.Unix()

	// modified
	t, err = time.Parse(parseDateFormat, p.ModifiedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其修改时间是出错：[%v]", name, err)
	}
	p.Modified = t.Unix()

	return p, nil
}

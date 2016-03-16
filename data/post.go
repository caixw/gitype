// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/issue9/logs"
	"gopkg.in/yaml.v2"
)

type Post struct {
	Slug           string  `yaml:"-"`        // 唯一名称
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
	Path           string  `yaml:"path"`     // 正文的文件名，相对于meta.yaml所在的目录
	Permalink      string  `yaml:"-"`        // 文章的唯一链接
}

// 加载所有的文章内容。
// dir data/posts目录。
func (d *Data) loadPosts(dir string) error {
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
	dir = filepath.Clean(dir)
	d.Posts = make([]*Post, 0, len(paths))
	for _, p := range paths {
		p = filepath.Clean(p)
		post, err := loadPost(dir, p, d.Config, d.Tags)
		if err != nil {
			logs.Error(err)
			continue
		}
		post.Permalink = path.Join(d.URLS.Root, d.URLS.Post, post.Slug+d.URLS.Suffix)

		d.Posts = append(d.Posts, post)
	}
	sort.Sort(posts(d.Posts))

	return nil
}

// 加载某一文章。
//
// postsDir 表示data/posts目录的绝对地址，必须经过filepath.Clean()处理；
// path 表示具体文章的meta.yaml文章，必须经过filepath.Clean()处理；
func loadPost(postsDir, path string, conf *Config, tags []*Tag) (*Post, error) {
	dir := filepath.Dir(path)                        // 获取路径部分
	slug := strings.TrimPrefix(dir, postsDir)        // 获取相对于data/posts的名称
	slug = strings.Trim(filepath.ToSlash(slug), "/") // 转换成/符号并去掉首尾的/字符

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := &Post{}
	if err := yaml.Unmarshal(data, p); err != nil {
		return nil, fmt.Errorf("[%v]解板yaml出错:%v\n", slug, err)
	}

	if len(p.Title) == 0 {
		return nil, fmt.Errorf("[%v]:文章标题不能为空\n", slug)
	}
	p.Slug = slug

	if p.Author == nil {
		p.Author = conf.Author
	}

	// content
	if len(p.Path) == 0 {
		return nil, fmt.Errorf("[%v]:未指定内容文件\n", slug)
	}
	data, err = ioutil.ReadFile(filepath.Join(dir, p.Path))
	if err != nil {
		return nil, fmt.Errorf("[%v]:读取文章内容出错：[%v]\n", slug, err)
	}
	p.Content = string(data)

	// tags
	ts := strings.Split(p.TagsString, ",")
	if len(ts) == 0 {
		return nil, fmt.Errorf("文章[%v]未指定任何关联标签信息\n", slug)
	}
	for _, tag := range tags {
		for _, tagName := range ts {
			if tag.Slug == tagName {
				p.Tags = append(p.Tags, tag)
				tag.Posts = append(tag.Posts, p)
				break
			} // end if
		} // end for ts
	} // end for tags

	// created
	t, err := time.Parse(parseDateFormat, p.CreatedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其创建时间是出错：[%v]\n", slug, err)
	}
	p.Created = t.Unix()

	// modified
	t, err = time.Parse(parseDateFormat, p.ModifiedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其修改时间是出错：[%v]\n", slug, err)
	}
	p.Modified = t.Unix()

	// 指定默认模板
	if len(p.Template) == 0 {
		p.Template = "post"
	}

	return p, nil
}

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

	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"gopkg.in/yaml.v2"
)

func (d *Data) postPath(path string) string {
	return filepath.Join(d.Root, "posts", path)
}

// 加载所有的文章内容。
// dir data/posts目录。
func (d *Data) loadPosts() error {
	dir := d.postPath("")
	paths := make([]string, 0, 100)

	// 遍历data/posts目录，查找所有的 meta.yaml 文章。
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
		post.Permalink = path.Join(d.Config.URLS.Root, d.Config.URLS.Post, post.Slug+d.Config.URLS.Suffix)

		d.Posts = append(d.Posts, post)
	}

	sort.SliceStable(d.Posts, func(i, j int) bool {
		switch {
		case d.Posts[i].Top && d.Posts[j].Top:
			return d.Posts[i].Created >= d.Posts[j].Created
		case d.Posts[i].Top:
			return false
		case d.Posts[j].Top:
			return true
		default:
			return d.Posts[i].Created >= d.Posts[j].Created
		}

	})

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
		return nil, fmt.Errorf("[%v]解板yaml出错:%v", slug, err)
	}

	if len(p.Title) == 0 {
		return nil, fmt.Errorf("[%v]:文章标题不能为空", slug)
	}
	p.Slug = slug

	if p.Author == nil {
		p.Author = conf.Author
	}

	// content
	if len(p.Path) == 0 {
		return nil, fmt.Errorf("[%v]:未指定内容文件", slug)
	}
	data, err = ioutil.ReadFile(filepath.Join(dir, p.Path))
	if err != nil {
		return nil, fmt.Errorf("[%v]:读取文章内容出错：[%v]", slug, err)
	}
	p.Content = string(data)

	// tags
	ts := strings.Split(p.TagsString, ",")
	if len(ts) == 0 {
		return nil, fmt.Errorf("文章[%v]未指定任何关联标签信息", slug)
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
	if len(p.Tags) == 0 {
		return nil, fmt.Errorf("文章[%v]未指定任何有效的关联标签信息", slug)
	}

	// created
	t, err := time.Parse(vars.DateFormat, p.CreatedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其创建时间是出错：[%v]", slug, err)
	}
	p.Created = t.Unix()

	// modified
	t, err = time.Parse(vars.DateFormat, p.ModifiedFormat)
	if err != nil {
		return nil, fmt.Errorf("[%v]:解析其修改时间是出错：[%v]", slug, err)
	}
	p.Modified = t.Unix()

	// 指定默认模板
	if len(p.Template) == 0 {
		p.Template = "post"
	}

	return p, nil
}

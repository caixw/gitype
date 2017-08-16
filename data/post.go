// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caixw/typing/vars"
	"gopkg.in/yaml.v2"
)

// 加载所有的文章内容。
func (d *Data) loadPosts() error {
	dir := d.path.PostsDir
	paths := make([]string, 0, 100)

	// 遍历 data/posts 目录，查找所有的 meta.yaml 文章。
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
			return err
		}
		post.Permalink = vars.PostURL(post.Slug)

		d.Posts = append(d.Posts, post)
	}

	return nil
}

// 加载某一文章。
//
// postsDir 表示 data/posts 目录的绝对地址，必须经过 filepath.Clean() 处理；
// path 表示具体文章的 meta.yaml 文章，必须经过 filepath.Clean() 处理。
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
		return nil, fmt.Errorf("%s 解板 yaml 出错:%v", slug, err)
	}

	if len(p.Title) == 0 {
		return nil, fmt.Errorf("%s:文章标题不能为空", slug)
	}
	p.Slug = slug

	if p.Author == nil {
		p.Author = conf.Author
	}

	// content
	if len(p.Path) == 0 {
		return nil, fmt.Errorf("%s:未指定内容文件", slug)
	}
	data, err = ioutil.ReadFile(filepath.Join(dir, p.Path))
	if err != nil {
		return nil, fmt.Errorf("%s:读取文章内容出错：%v", slug, err)
	}
	p.Content = string(data)

	// tags
	ts := strings.Split(p.TagsString, ",")
	if len(ts) == 0 {
		return nil, fmt.Errorf("文章 %s 未指定任何关联标签信息", slug)
	}
	for _, tag := range tags {
		for _, tagName := range ts {
			if tag.Slug == tagName {
				p.Tags = append(p.Tags, tag)
				tag.Posts = append(tag.Posts, p)
				break
			}
		}
	} // end for tags
	if len(p.Tags) == 0 {
		return nil, fmt.Errorf("文章 %s 未指定任何有效的关联标签信息", slug)
	}

	// keywords
	if len(p.Keywords) == 0 && len(p.Tags) > 0 {
		keywords := make([]string, 0, len(p.Tags))
		for _, v := range p.Tags {
			keywords = append(keywords, v.Title)
		}
		p.Keywords = strings.Join(keywords, ",")
	}

	// created
	p.Created, err = parseDate(p.CreatedFormat)
	if err != nil {
		return nil, fmt.Errorf("%s 解析其创建时间是出错：%v", slug, err)
	}

	// modified
	p.Modified, err = parseDate(p.ModifiedFormat)
	if err != nil {
		return nil, fmt.Errorf("%s 解析其修改时间是出错：%v", slug, err)
	}

	// 指定默认模板
	if len(p.Template) == 0 {
		p.Template = "post"
	}

	if len(p.Order) == 0 {
		p.Order = OrderDefault
	} else if p.Order != OrderDefault && p.Order != OrderLast && p.Order != OrderTop {
		return nil, fmt.Errorf("无效的 Order 值: %s", p.Order)
	}

	return p, nil
}

func sortPosts(posts []*Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		switch {
		case posts[i].Order == posts[j].Order:
			return posts[i].Created >= posts[j].Created
		case (posts[i].Order == OrderTop) || (posts[j].Order == OrderLast):
			return true
		case (posts[i].Order == OrderLast) || (posts[j].Order == OrderTop):
			return false
		default:
			return posts[i].Created >= posts[j].Created
		}
	})
}

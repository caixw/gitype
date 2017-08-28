// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caixw/typing/vars"
)

const (
	postMetaFilename    = "meta.yaml"
	postContentFilename = "content.html"
)

// 表示 Post.Order 的各类值
const (
	orderTop     = "top"     // 置顶
	orderLast    = "last"    // 放在尾部
	orderDefault = "default" // 默认情况
)

// Post 表示文章的信息
type Post struct {
	Slug           string `yaml:"-"`                  // 唯一名称
	Title          string `yaml:"title"`              // 标题
	Created        int64  `yaml:"-"`                  // 创建时间，用时间戳，而不是 time.Time，方便模板用户打印数据
	Modified       int64  `yaml:"-"`                  // 修改时间
	Tags           []*Tag `yaml:"-"`                  // 关联的标签
	Keywords       string `yaml:"keywords,omitempty"` // meta.keywords 标签的内容，如果为空，使用 tags
	Order          string `yaml:"order,omitempty"`    // 排序方式
	Summary        string `yaml:"summary"`            // 摘要，同时也作为 meta.description 的内容
	Content        string `yaml:"-"`                  // 内容
	CreatedFormat  string `yaml:"created"`            // 创建时间的字符串表示形式
	ModifiedFormat string `yaml:"modified"`           // 修改时间的字符串表示形式
	TagsString     string `yaml:"tags"`               // 关联标签的列表
	Path           string `yaml:"path"`               // 正文的文件名，相对于 meta.yaml 所在的目录
	Permalink      string `yaml:"-"`                  // 文章的唯一链接
	Outdated       string `yaml:"-"`                  // 已过时文章的提示信息，这是一个动态的值，不能提前计算

	// 以下内容不存在时，则会使用全局的默认选项
	Author   *Author `yaml:"author,omitempty"`   // 作者
	License  *Link   `yaml:"license,omitempty"`  // 版本信息
	Template string  `yaml:"template,omitempty"` // 使用的模板
}

func loadPosts(dir string) ([]*Post, error) {
	paths := make([]string, 0, 100)

	// 遍历 data/posts 目录，查找所有的 meta.yaml 文章。
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == postMetaFilename {
			paths = append(paths, path)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	// 开始加载文章的具体内容。
	dir = filepath.Clean(dir)
	posts := make([]*Post, 0, len(paths))
	for _, p := range paths {
		p = filepath.Clean(p)
		post, err := loadPost(dir, p)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func loadPost(postsDir, path string) (*Post, error) {
	dir := filepath.Dir(path)                 // 获取路径部分
	slug := strings.TrimPrefix(dir, postsDir) // 获取相对于 data/posts 的名称
	slug = strings.Trim(filepath.ToSlash(slug), "/")

	p := &Post{}
	if err := loadYamlFile(path, p); err != nil {
		return nil, err
	}
	p.Slug = slug

	// path
	if len(p.Path) == 0 {
		p.Path = postContentFilename
	}
	data, err := ioutil.ReadFile(filepath.Join(dir, p.Path))
	if err != nil {
		return nil, &FieldError{File: p.Slug, Message: err.Error(), Field: "path"}
	}
	p.Content = string(data)

	return p, nil
}

func (p *Post) sanitize() *FieldError {
	if len(p.Title) == 0 {
		return &FieldError{File: p.Slug, Message: "不能为空", Field: "title"}
	}

	// permalink
	p.Permalink = vars.PostURL(p.Slug)

	// content
	if len(p.Content) == 0 {
		return &FieldError{File: p.Slug, Message: "不能为空", Field: "content"}
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
	created, err := vars.ParseDate(p.CreatedFormat)
	if err != nil {
		return &FieldError{File: p.Slug, Message: err.Error(), Field: "created"}
	}
	p.Created = created
	p.CreatedFormat = ""

	// modified
	modified, err := vars.ParseDate(p.ModifiedFormat)
	if err != nil {
		return &FieldError{File: p.Slug, Message: err.Error(), Field: "modified"}
	}
	p.Modified = modified
	p.ModifiedFormat = ""

	// template
	if len(p.Template) == 0 {
		p.Template = "post"
	}

	// order
	if len(p.Order) == 0 {
		p.Order = orderDefault
	} else if p.Order != orderDefault && p.Order != orderLast && p.Order != orderTop {
		return &FieldError{File: p.Slug, Message: "无效的值", Field: "order"}
	}

	return nil
}

func sortPosts(posts []*Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		switch {
		case posts[i].Order == posts[j].Order:
			return posts[i].Created >= posts[j].Created
		case (posts[i].Order == orderTop) || (posts[j].Order == orderLast):
			return true
		case (posts[i].Order == orderLast) || (posts[j].Order == orderTop):
			return false
		default:
			return posts[i].Created >= posts[j].Created
		}
	})
}

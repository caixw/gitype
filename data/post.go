// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
	"github.com/issue9/utils"
)

// 文章是否过时的比较方式
const (
	outdatedTypeCreated  = "created"
	outdatedTypeModified = "modified"
	outdatedTypeNone     = "none"
	outdatedTypeCustom   = "custom"
)

// 表示 Post.State 的各类值
const (
	stateTop     = "top"     // 置顶
	stateLast    = "last"    // 放在尾部
	stateDefault = "default" // 默认值
	stateDraft   = "draft"   // 表示为草稿，不会加载此条数据
)

// Post 表示文章的信息
type Post struct {
	Slug      string    `yaml:"-"`                  // 唯一名称
	Permalink string    `yaml:"created"`            // 文章的唯一链接，同时当作 created 的原始值
	Title     string    `yaml:"title"`              // 标题
	HTMLTitle string    `yaml:"modified"`           // 网页标题，同时当作 modified 的原始值
	Created   time.Time `yaml:"-"`                  // 创建时间
	Modified  time.Time `yaml:"-"`                  // 修改时间
	Summary   string    `yaml:"summary"`            // 摘要，同时也作为 meta.description 的内容
	Content   string    `yaml:"outdated,omitempty"` // 内容，同时也作为 outdated 的内容

	// 关联的标签列表，以半角逗号分隔的字符串，
	// 标签名为各个标签的 slug 值，可以保证其唯一。
	// 最终会被解析到 Tags 中，TagString 会被废弃。
	TagsString string `yaml:"tags"`
	Tags       []*Tag `yaml:"-"`

	// Outdated 用户记录文章的一个过时情况，可以由以下几种值构成：
	// - created 表示该篇文章以创建时间来计算其是否已经过时，该值也是默认值；
	// - modified 表示该文章以其修改时间来计算其是否已经过时；
	// - none 表示该文章永远不会过时；
	// - 其它任意非空值，表示直接以该字符串当作过时信息展示给用语.
	// yaml 内容由 Content 字段读入。
	Outdated *Outdated `yaml:"-"`

	// State 表示文章的状态，有以下四种值：
	// - top 表示文章被置顶；
	// - last 表示文章会被放置在最后；
	// - draft 表示这是一篇草稿，并不会被加地到内存中；
	// - default 表示默认情况，也可以为空，按默认的方式进行处理。
	State string `yaml:"state,omitempty"`

	// 以下内容不存在时，则会使用全局的默认选项
	Author   *Author `yaml:"author,omitempty"`
	License  *Link   `yaml:"license,omitempty"`
	Template string  `yaml:"template,omitempty"`
	Keywords string  `yaml:"keywords,omitempty"`

	// 用于搜索的副本内容，会全部转换成小写
	SearchTitle   string
	SearchContent string
}

// Outdated 表示每一篇文章的过时情况
type Outdated struct {
	Type string
	Date time.Time
	Days int

	Content string // 自定义的提示内容
}

func loadPosts(path *path.Path) ([]*Post, error) {
	dir := path.PostsDir
	slugs := make([]string, 0, 100)

	// 遍历 data/posts 目录，查找所有的文章。
	walk := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		postsDir := filepath.Clean(path.PostsDir)
		slug := strings.TrimPrefix(p, postsDir) // 获取相对于 data/posts 的名称
		slug = strings.Trim(filepath.ToSlash(slug), "/")

		if utils.FileExists(path.PostContentPath(slug)) &&
			utils.FileExists(path.PostMetaPath(slug)) {
			slugs = append(slugs, slug)
		}
		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	// 开始加载文章的具体内容。
	posts := make([]*Post, 0, len(slugs))
	for _, slug := range slugs {
		post, err := loadPost(path, slug)
		if err != nil {
			return nil, err
		}

		if post.State != stateDraft {
			posts = append(posts, post)
		}
	}

	if err := checkPostsDup(posts); err != nil {
		return nil, err
	}

	sortPosts(posts)

	return posts, nil
}

func loadPost(path *path.Path, slug string) (*Post, error) {
	post := &Post{}
	if err := helper.LoadYAMLFile(path.PostMetaPath(slug), post); err != nil {
		return nil, err
	}
	if post.State == stateDraft {
		return post, nil
	}

	// slug
	post.Slug = slug

	// created
	// permalink 还用作其它功能，需要首先解析其值
	created, err := time.Parse(vars.DateFormat, post.Permalink)
	if err != nil {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: err.Error(), Field: "created"}
	}
	post.Created = created

	// permalink
	post.Permalink = vars.PostURL(post.Slug)

	// modified
	// HTMLTitle 还用作其它功能，需要首先解析其值
	modified, err := time.Parse(vars.DateFormat, post.HTMLTitle)
	if err != nil {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: err.Error(), Field: "modified"}
	}
	post.Modified = modified
	post.HTMLTitle = ""

	// outdated，依赖 modified 和 created
	// Content 还用作其它功能，需要首先解析其值
	switch post.Content {
	case outdatedTypeCreated, "":
		post.Outdated = &Outdated{
			Type: outdatedTypeCreated,
			Date: created,
		}
	case outdatedTypeModified:
		post.Outdated = &Outdated{
			Type: outdatedTypeModified,
			Date: modified,
		}
	case outdatedTypeNone:
		post.Outdated = nil
	default:
		post.Outdated = &Outdated{
			Type:    outdatedTypeCustom,
			Content: post.Content,
		}
	}
	post.Content = ""

	// 加载内容
	data, err := ioutil.ReadFile(path.PostContentPath(slug))
	if err != nil {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: err.Error(), Field: "path"}
	}
	if len(data) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "content"}
	}
	post.Content = string(data)

	if len(post.Title) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "title"}
	}

	if len(post.TagsString) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "tags"}
	}

	// keywords
	if len(post.Keywords) == 0 {
		post.Keywords = post.TagsString
	}

	// template
	if len(post.Template) == 0 {
		post.Template = vars.PagePost
	}

	// state
	if len(post.State) == 0 {
		post.State = stateDefault
	} else if post.State != stateDefault &&
		post.State != stateLast &&
		post.State != stateTop {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "无效的值", Field: "order"}
	}

	post.SearchContent = strings.ToLower(post.Content)
	post.SearchTitle = strings.ToLower(post.Title)

	return post, nil
}

// 检测是否存在同名的文章
func checkPostsDup(posts []*Post) error {
	count := func(slug string) (cnt int) {
		for _, post := range posts {
			if post.Slug == slug {
				cnt++
			}
		}
		return cnt
	}

	for _, post := range posts {
		if count(post.Slug) > 1 {
			return errors.New("存在同名的文章：" + post.Slug)
		}
	}

	return nil
}

// 对文章进行排序，需保证 created 已经被初始化
func sortPosts(posts []*Post) {
	sort.SliceStable(posts, func(i, j int) bool {
		switch {
		case (posts[i].State == stateTop) || (posts[j].State == stateLast):
			return true
		case (posts[i].State == stateLast) || (posts[j].State == stateTop):
			return false
		default:
			return posts[i].Created.After(posts[j].Created)
		}
	})
}

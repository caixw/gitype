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

	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/path"
	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// 文章是否过时的比较方式
const (
	outdatedTypeCreated  = "created"  // 以创建时间作为对比
	outdatedTypeModified = "modified" // 以修改时间作为对比
)

// 表示 Post.Order 的各类值
const (
	orderTop     = "top"     // 置顶
	orderLast    = "last"    // 放在尾部
	orderDefault = "default" // 默认情况
)

// 描述过时文章的提示信息。
//
// 理论上把有关 outdatedConfig 的信息，直接在模板中对文章的创建时间戳进行比较，
// 是比通过配置来比较会更加方便，也不会更任何的后期工作量。之所以把这个功能放在后端，
// 而不是模板层面，是因为觉得模板应该只负责展示页面，而不是用于处理逻辑内容。
type outdatedConfig struct {
	Type     string        `yaml:"type"`     // 比较的类型，创建时间或是修改时间
	Duration time.Duration `yaml:"duration"` // 超时的时间，可以使用 time.Duration 的字符串值
	Content  string        `yaml:"content"`  // 提示的内容，普通文字，不能为 html
}

// Post 表示文章的信息
type Post struct {
	Slug       string    `yaml:"-"`               // 唯一名称
	Title      string    `yaml:"title"`           // 标题
	Created    time.Time `yaml:"-"`               // 创建时间
	Modified   time.Time `yaml:"-"`               // 修改时间
	Tags       []*Tag    `yaml:"-"`               // 关联的标签
	Summary    string    `yaml:"summary"`         // 摘要，同时也作为 meta.description 的内容
	Content    string    `yaml:"-"`               // 内容
	TagsString string    `yaml:"tags"`            // 关联标签的列表
	Permalink  string    `yaml:"created"`         // 文章的唯一链接，同时当作 created 的原始值
	Outdated   string    `yaml:"modified"`        // 已过时文章的提示信息，同时当作 modified 的原始值
	Order      string    `yaml:"order,omitempty"` // 排序方式
	Draft      bool      `yaml:"draft,omitempty"` // 是否为草稿，为 true，则不会加载该条数据

	// 以下内容不存在时，则会使用全局的默认选项
	Author   *Author `yaml:"author,omitempty"`   // 作者
	License  *Link   `yaml:"license,omitempty"`  // 版本信息
	Template string  `yaml:"template,omitempty"` // 使用的模板
	Keywords string  `yaml:"keywords,omitempty"` // meta.keywords 标签的内容，如果为空，使用 tags
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

		if !post.Draft {
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
	if post.Draft {
		return post, nil
	}

	// slug
	post.Slug = slug

	// 加载内容
	data, err := ioutil.ReadFile(path.PostContentPath(slug))
	if err != nil {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: err.Error(), Field: "path"}
	}
	if len(data) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "content"}
	}
	post.Content = string(data)

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
	// outdated 还用作其它功能，需要首先解析其值
	modified, err := time.Parse(vars.DateFormat, post.Outdated)
	if err != nil {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: err.Error(), Field: "modified"}
	}
	post.Modified = modified
	post.Outdated = ""

	if len(post.Title) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "title"}
	}

	if len(post.TagsString) == 0 {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "不能为空", Field: "tags"}
	}

	// keywords
	if len(post.Keywords) == 0 && len(post.Tags) > 0 {
		keywords := make([]string, 0, len(post.Tags))
		for _, v := range post.Tags {
			keywords = append(keywords, v.Title)
		}
		post.Keywords = strings.Join(keywords, ",")
	}

	// template
	if len(post.Template) == 0 {
		post.Template = vars.PostTemplateName
	}

	// order
	if len(post.Order) == 0 {
		post.Order = orderDefault
	} else if post.Order != orderDefault && post.Order != orderLast && post.Order != orderTop {
		return nil, &helper.FieldError{File: path.PostMetaPath(slug), Message: "无效的值", Field: "order"}
	}

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
		case (posts[i].Order == orderTop) || (posts[j].Order == orderLast):
			return true
		case (posts[i].Order == orderLast) || (posts[j].Order == orderTop):
			return false
		default:
			return posts[i].Created.After(posts[j].Created)
		}
	})
}

func (o *outdatedConfig) sanitize() *helper.FieldError {
	if o.Type != outdatedTypeCreated && o.Type != outdatedTypeModified {
		return &helper.FieldError{Message: "无效的值", Field: "outdated.type"}
	}

	if len(o.Content) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "outdated.content"}
	}

	return nil
}

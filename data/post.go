// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"sort"
	"strings"
	"time"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
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
	Slug      string    // 唯一名称
	Permalink string    // 文章的唯一链接，同时当作 created 的原始值
	Title     string    // 标题
	HTMLTitle string    // 网页标题，同时当作 modified 的原始值
	Created   time.Time // 创建时间
	Modified  time.Time // 修改时间
	Summary   string    // 摘要，同时也作为 meta.description 的内容
	Content   string    // 内容，同时也作为 outdated 的内容
	Tags      []*Tag
	Outdated  *Outdated
	State     string

	// 以下内容不存在时，则会使用全局的默认选项
	Author   *Author
	License  *Link
	Template string
	Keywords string
}

// Outdated 表示每一篇文章的过时情况
type Outdated struct {
	Type string
	Date time.Time
	Days int

	Content string // 自定义的提示内容
}

func loadPosts(path *path.Path, tags []*Tag, conf *loader.Config) ([]*Post, error) {
	ps, err := loader.LoadPosts(path)
	if err != nil {
		return nil, err
	}

	// 开始加载文章的具体内容。
	posts := make([]*Post, 0, len(ps))
	for _, p := range ps {
		if p.State == loader.StateDraft { // 草稿不收录
			continue
		}

		post := &Post{
			Slug:      p.Slug,
			Permalink: vars.PostURL(p.Slug),
			Title:     p.Title,
			Created:   p.Created,
			Modified:  p.Modified,
			Summary:   p.Summary,
			Content:   p.Content,
			State:     p.State,
			HTMLTitle: helper.ReplaceContent(conf.Pages[vars.PagePost].Title, p.Title),

			Author:   p.Author,
			License:  p.License,
			Template: p.Template,
			Keywords: p.Keywords,
		}

		switch p.Outdated {
		case loader.OutdatedTypeCreated, "":
			post.Outdated = &Outdated{
				Type: outdatedTypeCreated,
				Date: post.Created,
			}
		case loader.OutdatedTypeModified:
			post.Outdated = &Outdated{
				Type: outdatedTypeModified,
				Date: post.Modified,
			}
		case loader.OutdatedTypeNone:
			post.Outdated = nil
		default:
			post.Outdated = &Outdated{
				Type:    outdatedTypeCustom,
				Content: post.Content,
			}
		}

		if post.Author == nil {
			post.Author = conf.Author
		}

		if post.License == nil {
			post.License = conf.License
		}

		if err := attachPostTag(path, post, tags, p.Tags); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	sortPosts(posts)

	return posts, nil
}

// 关联文章与标签的相关信息
func attachPostTag(p *path.Path, post *Post, tags []*Tag, tagString string) *helper.FieldError {
	ts := strings.Split(tagString, ",")
	for _, tag := range tags {
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
		return &helper.FieldError{File: p.PostMetaPath(post.Slug), Message: "未指定任何关联标签信息", Field: "tags"}
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

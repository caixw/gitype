// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"strconv"
	"strings"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// 文章的详细内容
type Post struct {
	ID           int64
	Name         string
	Title        string
	Summary      string
	Content      string
	Author       string // 作者名称
	Created      int64  // 创建时间
	Modified     int64  // 修改时间
	AllowComment bool   // 是否允许评论
}

func (p *Post) CommentsSize() int {
	if size, found := stat.Posts[p.ID]; found {
		return size
	}

	c := &models.Comment{PostID: p.ID, State: models.CommentStateApproved}
	size, err := db.Count(c)
	if err != nil {
		logs.Error("front.Post.CommentsSize:", err)
		return 0
	}
	stat.Posts[p.ID] = size
	return size
}

// 返回文章的摘要或是具体内容。
func (p *Post) Entry() string {
	if len(p.Summary) > 0 {
		return p.Summary
	}
	return p.Content
}

// 返回文章的链接，相对URL
func (p *Post) Permalink() string {
	if len(p.Name) > 0 {
		return opt.PostURL(p.Name)
	}
	return opt.PostURL(strconv.FormatInt(p.ID, 10))
}

// 与该文章相关的关键字，即该文章所有标签的组成的字符串
func (p *Post) Keywords() string {
	tags := p.Tags()
	if len(tags) == 0 {
		return ""
	}

	keys := make([]string, 0, len(tags))
	for _, val := range tags {
		keys = append(keys, val.Title)
	}
	return strings.Join(keys, ",")
}

// 获取与当前文章相关的标签。
func (p *Post) Tags() []*Tag {
	sql := `SELECT t.{name} AS {Name}, t.{title} AS {Title} FROM #relationships AS r
	LEFT JOIN #tags AS t on t.{id}=r.{tagID}
	WHERE r.{postID}=?`

	rows, err := db.Query(true, sql, p.ID)
	if err != nil {
		logs.Error("front.Post.Tags:", err)
		return nil
	}
	defer rows.Close()

	tags := make([]*Tag, 0, 5)
	if _, err = fetch.Obj(&tags, rows); err != nil {
		logs.Error("front.Post.Tags:", err)
		return nil
	}
	return tags
}

// 返回文章的评论信息。
func (p *Post) Comments() []*Comment {
	sql := `SELECT {id} AS {ID}, {created} AS {Created}, {agent} AS {Agent}, {content} AS {Content},
	{isAdmin} AS {IsAdmin}, {authorName} AS {AuthorName}, {authorURL} AS {AuthorURL}, {postID} AS {PostID}
	FROM #comments
	WHERE {postID}=? AND {state}=?
	ORDER BY {created} `
	if opt.CommentOrder == app.CommentOrderDesc {
		sql += `DESC `
	}

	rows, err := db.Query(true, sql, p.ID, models.CommentStateApproved)
	if err != nil {
		logs.Error("front.Post.Comment:", err)
		return nil
	}
	defer rows.Close()

	comments := make([]*Comment, 0, opt.PageSize)
	if _, err := fetch.Obj(&comments, rows); err != nil {
		logs.Error("front.Post.Comment:", err)
		return nil
	}
	return comments
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"strconv"

	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// 一个锚点的表示形式。
type Anchor interface {
	Link() string // 链接
	Text() string // 字面文字
	Ext() string  // 扩展内容，比如title,alt等，根据链接来确定
}

// 页面的基本信息
type PageInfo struct {
	Title       string // 网页的title值
	SiteName    string // 网站名称
	SecondTitle string // 副标题
	Canonical   string // 当前页的唯一链接
	Keywords    string // meta.keywords的值
	Description string // meta.description的值
	AppVersion  string // 当前程序的版本号
	GoVersion   string // 编译的go版本号
	Author      string // 作者名称

	PostSize    int      // 文章数量
	CommentSize int      // 评论数量
	Tags        []Anchor // 标签列表
	Topics      []Anchor // 最新评论的10条内容
	Hots        []Anchor // 评论最多的10条内容
}

type Tag struct {
	Name        string
	Title       string
	Description string
}

func (t *Tag) Link() string {
	return opt.SiteURL + "/tags/" + t.Name
}
func (t *Tag) Text() string {
	return t.Title
}

func (t *Tag) Ext() string {
	return t.Description
}

// 文章的详细内容
type Post struct {
	ID           int64
	Name         string
	Title        string
	Content      string
	Author       string
	Comments     int   // 评论数量
	Created      int64 // 创建时间
	Modified     int64 // 修改时间
	AllowComment bool  // 是否允许评论
}

func (p *Post) Tags() []*Tag {
	sql := `SELECT t.{name} AS Name, t.{title} AS Text FROM #relationships AS r
	 LEFT JOIN #tags AS t on t.{id}=r.{tagID}
	 WHERE r.{postID}=?`

	rows, err := db.Query(true, sql, p.ID)
	if err != nil {
		logs.Error("themes.Post.Tags:", err)
		return nil
	}
	defer rows.Close()

	tags := make([]*Tag, 0, 5)
	if _, err = fetch.Obj(&tags, rows); err != nil {
		logs.Error("themes.Post.Tags:", err)
		return nil
	}
	return tags
}

func (p *Post) Permalink() string {
	if len(p.Name) > 0 {
		return opt.SiteURL + "/posts/" + p.Name + opt.Suffix
	}

	return opt.SiteURL + "/posts/" + strconv.FormatInt(p.ID, 10) + opt.Suffix
}

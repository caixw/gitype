// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

////////////////////////// comments

const (
	CommentStateAll      = iota // 表示所有以下的状态。
	CommentStateWaiting         // 等待审核
	CommentStateSpam            // 垃圾评论
	CommentStateApproved        // 通过验证
)

type Comment struct {
	ID      int64  `orm:"name(id);ai"`
	Parent  int64  `orm:"name(parent)"`          // 子评论的话，这此为其上一级评论的id
	Created int64  `orm:"name(created)"`         // 记录创建的时间
	PostID  int64  `orm:"name(postID)"`          // 被评论的文章id
	State   int    `orm:"name(state)"`           // 此条记录的状态
	IP      string `orm:"name(ip);len(50)"`      // 评论者的ip
	Agent   string `orm:"name(agent);len(200)"`  // 评论者的agent
	Content string `orm:"name(content);len(-1)"` // 评论内容

	IsAdmin     bool   `orm:"name(isAdmin)"`              // 网站的管理员评论
	AuthorName  string `orm:"name(authorName);len(20)"`   // 评论用户的名称
	AuthorEmail string `orm:"name(authorEmail);len(200)"` // 作者邮件地址
	AuthorURL   string `orm:"name(authorURL);len(200)"`   // 作者站点
}

func (c *Comment) Meta() string {
	return `name(comments)`
}

////////////////////////// post

const (
	PostStateAll       = iota // 表示所有状态
	PostStatePublished        // 已经发布
	PostStateDraft            // 草稿
)

// 文章内容
type Post struct {
	ID       int64  `orm:"name(id);ai"`
	Name     string `orm:"name(name);len(200);index(idx_name)"` // 唯一名称，为空表示使用id
	Title    string `orm:"name(title);len(200)"`                // 标题
	Summary  string `orm:"name(summary);len(5000)"`             // 内容摘要
	Content  string `orm:"name(content);len(-1)"`               // 实际内容
	State    int    `orm:"name(state)"`                         // 状态
	Order    int    `orm:"name(order)"`                         // 排序
	Template string `orm:"name(template);len(50)"`              // 使用的模板

	Created  int64 `orm:"name(created)"`  // 创建时间
	Modified int64 `orm:"name(modified)"` // 最后次修改时间

	AllowPing    bool `orm:"name(allowPing)"`
	AllowComment bool `orm:"name(allowComment)"`

	//License string `orm:"name(license)"`
}

func (p *Post) Meta() string {
	return `name(posts)`
}

////////////////////// option

// 系统设置项。
type Option struct {
	Key   string `orm:"name(key);len(20);pk"` // 该设置项的唯一名称
	Value string `orm:"name(value);len(-1)"`  // 该设置项的值
	Group string `orm:"name(group);len(20)"`  // 对该设置项的分组。
}

func (opt *Option) Meta() string {
	return `name(options)`
}

///////////////////// relationship

type Relationship struct {
	PostID int64 `orm:"name(postID);pk"`
	TagID  int64 `orm:"name(tagID);pk"`
}

func (r *Relationship) Meta() string {
	return `name(relationships)`
}

/////////////////////// tag

type Tag struct {
	ID          int64  `orm:"name(id);ai" json:"id"`
	Name        string `orm:"name(name);unique(u_name);len(50);" json:"name,omitempty"` // 唯一名称
	Title       string `orm:"name(title);unique(u_title);len(50)" json:"title"`         // 名称
	Description string `orm:"name(description);len(-1)" json:"description"`             // 详细描述，可以用html
}

func (t *Tag) Meta() string {
	return `name(tags)`
}

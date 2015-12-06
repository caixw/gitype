// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

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

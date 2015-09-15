// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"time"

	"github.com/caixw/typing/core"
	"github.com/issue9/logs"
)

const (
	commentStateAll      = iota // 表示所有以下的状态。
	commentStateWaiting         // 等待审核
	commentStateSpam            // 垃圾评论
	commentStateApproved        // 通过验证
)

type comment struct {
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

func (c *comment) Meta() string {
	return `name(comments)`
}

// @api get /admin/api/comments 获取所有评论内容
// @apiQuery page int
// @apiQuery size int
// @apiQuery state int
// @apiGroup admin
//
// @apiSuccess 200 ok
// @apiParam count int 符合条件(去除page和size条件)的所有评论数量
// @apiParam comments array 当前页的评论
func getComments(w http.ResponseWriter, r *http.Request) {
	var page, size, state int
	var ok bool
	if state, ok = core.QueryInt(w, r, "state", commentStateAll); !ok {
		return
	}

	sql := db.SQL().Table("#comments")
	if state != commentStateAll {
		sql.And("{state}=?", state)
	}
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("getComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	if size, ok = core.QueryInt(w, r, "size", opt.PageSize); !ok {
		return
	}
	sql.Limit(size, page*size)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("getComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "comments": maps}, nil)
}

// @api put /admin/api/comments/{id} 修改评论，只能修改管理员发布的评论
// @apiParam id int 需要修改的评论id
// @apiGroup admin
//
// @apiRequest json
// @apiParam content string 新的评论内容
// @apiExample json
// { "content", "content..." }
//
// @apiSuccess 200 ok
func putComment(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	c := &comment{ID: id}
	cnt, err := db.Count(c)
	if err != nil {
		logs.Error("putComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if cnt == 0 {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	ct := &struct {
		Content string `json:"content"`
	}{}

	if !core.ReadJSON(w, r, ct) {
		return
	}

	c.Content = ct.Content

	if _, err = db.Update(c); err != nil {
		logs.Error("putComment", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api post /admin/api/comments/{id}/waiting 将评论的状态改为waiting
// @apiGroup admin
//
// @apiSuccess 204 no content
func setCommentWaiting(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, commentStateWaiting)
}

// @api post /admin/api/comments/{id}/spam 将评论的状态改为spam
// @apiGroup admin
//
// @apiSuccess 204 no content
func setCommentSpam(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, commentStateSpam)
}

// @api post /admin/api/comments/{id}/approved 将评论的状态改为approved
// @apiGroup admin
//
// @apiSuccess 204 no content
func setCommentApproved(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, commentStateApproved)
}

func setCommentState(w http.ResponseWriter, r *http.Request, state int) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	c := &comment{ID: id, State: state}
	if _, err := db.Update(c); err != nil {
		logs.Error("setCommentState:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api post /admin/api/comments 提交新评论
// @apiGroup admin
//
// @apiRequest json
// @apiParam parent int 评论的父级内容
// @apiParam postID int 评论的文章
// @apiParam content string 评论的内容
//
// @apiSuccess 201 created
func adminPostComment(w http.ResponseWriter, r *http.Request) {
	c := &struct {
		Parent  int64  `json:"parent"`
		PostID  int64  `json:"postID"`
		Content string `json:"content"`
	}{}

	if !core.ReadJSON(w, r, c) {
		return
	}

	comm := &comment{
		Parent:      c.Parent,
		PostID:      c.PostID,
		Content:     c.Content,
		State:       commentStateApproved,
		IP:          "",
		Agent:       "",
		Created:     time.Now().Unix(),
		IsAdmin:     true,
		AuthorURL:   opt.SiteURL,
		AuthorName:  opt.ScreenName,
		AuthorEmail: opt.Email,
	}
	if _, err := db.Insert(comm); err != nil {
		logs.Error("adminPostComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, nil, nil)
}

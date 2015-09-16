// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
)

// @api get /admin/api/comments 获取所有评论内容
// @apiQuery page int
// @apiQuery size int
// @apiQuery state int
// @apiGroup admin
//
// @apiSuccess 200 ok
// @apiParam count int 符合条件(去除page和size条件)的所有评论数量
// @apiParam comments array 当前页的评论
func adminGetComments(w http.ResponseWriter, r *http.Request) {
	var page, size, state int
	var ok bool
	if state, ok = core.QueryInt(w, r, "state", models.CommentStateAll); !ok {
		return
	}

	sql := db.SQL().Table("#comments")
	if state != models.CommentStateAll {
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
func adminPutComment(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	c := &models.Comment{ID: id}
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
func adminSetCommentWaiting(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, models.CommentStateWaiting)
}

// @api post /admin/api/comments/{id}/spam 将评论的状态改为spam
// @apiGroup admin
//
// @apiSuccess 204 no content
func adminSetCommentSpam(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, models.CommentStateSpam)
}

// @api post /admin/api/comments/{id}/approved 将评论的状态改为approved
// @apiGroup admin
//
// @apiSuccess 204 no content
func adminSetCommentApproved(w http.ResponseWriter, r *http.Request) {
	setCommentState(w, r, models.CommentStateApproved)
}

func setCommentState(w http.ResponseWriter, r *http.Request, state int) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	c := &models.Comment{ID: id, State: state}
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

	comm := &models.Comment{
		Parent:      c.Parent,
		PostID:      c.PostID,
		Content:     c.Content,
		State:       models.CommentStateApproved,
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

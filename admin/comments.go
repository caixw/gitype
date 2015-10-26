// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api delete /admin/api/comments/{id} 删除某条评论
// @apiParam id int 评论的id值
// @apiGroup admin
//
// @apiSuccess 204 no content
func adminDeleteComment(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	c := &models.Comment{ID: id}
	if _, err := db.Delete(c); err != nil {
		logs.Error("adminDeleteComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api get /admin/api/comments 获取所有评论内容
// @apiQuery page  int 显示第page页的内容，基数0;
// @apiQuery size  int 每页显示的数量；
// @apiQuery state int 仅显示状态值为state的记录；
// @apiGroup admin
//
// @apiSuccess 200 OK
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
		logs.Error("adminGetComments:", err)
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
	maps, err := sql.SelectMapString(true, "*")
	if err != nil {
		logs.Error("adminGetComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "comments": maps}, nil)
}

// @api get /admin/api/comments/count 获取各种状态下的评论数量
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam all      int 评论的总量
// @apiParam waiting  int 等待审核的评论数量
// @apiParam spam     int 垃圾评论的数量
// @apiParam approved int 通过审核的评论数量
func adminGetCommentsCount(w http.ResponseWriter, r *http.Request) {
	sql := "SELECT {state}, count(*) AS cnt FROM #comments GROUP BY {state}"
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("adminGetCommentsCount:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	defer rows.Close()

	maps, err := fetch.MapString(false, rows)
	if err != nil {
		logs.Error("adminGetCommentsCount:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := map[string]int{
		"all":      0,
		"waiting":  0,
		"spam":     0,
		"approved": 0,
	}
	count := 0
	for _, v := range maps {
		num, err := strconv.Atoi(v["cnt"])
		if err != nil {
			logs.Error("adminGetCommentsCount:", err)
			core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
			return
		}
		count += num
		switch v["state"] {
		case "1":
			data["waiting"] = num
		case "2":
			data["spam"] = num
		case "3":
			data["approved"] = num
		default:
			logs.Error("adminGetCommentsCount: 未知的评论状态:", v["state"])
		}
	}
	data["all"] = count // 所有评论的数量
	core.RenderJSON(w, http.StatusOK, data, nil)
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

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/is"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api get /api/posts/{id} 获取某一文章的详细内容
// @apiGroup front
//
// @apiSuccess 200 ok
// @apiParam id           int    id值
// @apiParam type         int    文章类型
// @apiParam name         string 唯一名称，可以为空
// @apiParam title        string 标题
// @apiParam content      string 文章内容
// @apiParam created      int    创建时间
// @apiParam modified     int    修改时间
// @apiParam template     string 所使用的模板
// @apiParam allowPing    bool   允许ping
// @apiParam allowComment bool   允许评论
// @apiParam tags         array  文章关联的标签
func frontGetPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("getPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	tags, err := getPostTags(id)
	if err != nil {
		logs.Error("getPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	obj := &struct {
		ID           int64   `json:"id"`
		Type         int     `json:"type"`
		Name         string  `json:"name"`
		Title        string  `json:"title"`
		Content      string  `json:"content"`
		Created      int64   `json:"created"`
		Modified     int64   `json:"modified"`
		Template     string  `json:"template"`
		AllowPing    bool    `json:"AllowPing"`
		AllowComment bool    `json:"AllowComment"`
		Tags         []int64 `json:"tags"`
	}{
		ID:           p.ID,
		Type:         p.Type,
		Name:         p.Name,
		Title:        p.Title,
		Content:      p.Content,
		Created:      p.Created,
		Modified:     p.Modified,
		Template:     p.Template,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Tags:         tags,
	}
	core.RenderJSON(w, http.StatusOK, obj, nil)
}

// 获取与某post相关联的标签
func getPostTags(postID int64) ([]int64, error) {
	sql := `SELECT rs.{tagID} FROM #relationships AS rs
	LEFT JOIN #tags AS t ON t.{id}=rs.{tagID}
	WHERE rs.{postID}=?`
	rows, err := db.Query(true, sql, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "tagID", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, v := range maps {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}
	return ret, nil
}

// @api get /api/posts 获取前端可访问的文章列表
// @apiQuery page int
// @apiQuery size int
// @apiGroup front
//
// @apiSuccess 200 OK
func frontGetPosts(w http.ResponseWriter, r *http.Request) {
	sql := db.Where("{state}=?", models.PostStatePublished).Table("#posts")
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("frontGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	var page, size int
	var ok bool
	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	if size, ok = core.QueryInt(w, r, "size", opt.PageSize); !ok {
		return
	}
	sql.Limit(size, page*size)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("frontGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "posts": maps}, nil)
}

// @api get /api/posts/{id}/comments
// @apiQuery page  int 页码
// @apiQuery size  int 每页显示的数量
// @apiQuery order int 排序方式
// @apiGroup front
//
// @apiSuccess 200 OK
// @apiParam count    int   当前文章的所有评论数量
// @apiParam comments array 当前页的评论
func frontGetPostComments(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("frontGetPostComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	sql := db.Where("{postID}=?", id).
		And("{state}=?", models.CommentStateApproved).
		Table("#comments")

	var order, page, size int
	if order, ok = core.QueryInt(w, r, "order", core.CommentOrderUndefined); !ok {
		return
	}
	switch order {
	case core.CommentOrderAsc:
		sql.Asc("{order}")
	case core.CommentOrderDesc:
		sql.Desc("{order}")
	case core.CommentOrderUndefined:
	default:
		errs := &core.ErrorResult{Message: "格式错误"}
		errs.Detail["order"] = "取值错误，只能是0,1,2"
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("frontGetPostComments:", err)
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

// @api post /api/posts/{id}/comments 提交新评论
// @apiGroup front
//
// @apiRequest json
// @apiParam parent      int    评论的父级内容
// @apiParam postID      int    评论的文章
// @apiParam content     string 评论的内容
// @apiParam authorName  string 评论的作者
// @apiParam authorURL   string 评论作者的网站地址，可为空
// @apiParam authorEmail string 评论作者的邮箱
//
// @apiSuccess 201 created
func frontPostPostComment(w http.ResponseWriter, r *http.Request) {
	c := &struct {
		Parent      int64  `json:"parent"`
		PostID      int64  `json:"postID"`
		Content     string `json:"content"`
		AuthorName  string `json:"authorName"`
		AuthorURL   string `json:"authorURL"`
		AuthorEmail string `json:"authorEmail"`
	}{}

	if !core.ReadJSON(w, r, c) {
		return
	}

	// 判断文章状态
	if c.PostID <= 0 {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	p := &models.Post{ID: c.PostID}
	if err := db.Select(p); err != nil {
		logs.Error("forntPostPostComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if (len(p.Title) == 0 && len(p.Content) == 0) || p.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}
	if !p.AllowComment {
		core.RenderJSON(w, http.StatusMethodNotAllowed, nil, nil)
		return
	}

	// 判断提交数据的状态
	errs := &core.ErrorResult{}
	if c.Parent < 0 {
		errs.Detail["parent"] = "无效的parent"
	}
	if len(c.Content) == 0 {
		errs.Detail["content"] = "content不能为空"
	}
	if len(c.AuthorURL) > 0 && !is.URL(c.AuthorURL) {
		errs.Detail["authorURL"] = "无效的authorURL"
	}
	if !is.Email(c.AuthorEmail) {
		errs.Detail["authorEmail"] = "无效的authorEmail"
	}
	if len(c.AuthorName) == 0 {
		errs.Detail["authorName"] = "authorName不能为空"
	}

	c.AuthorName = html.EscapeString(c.AuthorName)

	// url只提取其host部分，其余的都去掉
	u, err := url.Parse(c.AuthorURL)
	if err != nil {
		logs.Error("postComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	c.AuthorURL = u.Scheme + ":" + u.Host

	c.Content = html.EscapeString(c.Content)
	c.Content = strings.Replace(c.Content, "\n", "<br />", -1)

	comm := &models.Comment{
		PostID:      c.PostID,
		Parent:      c.Parent,
		AuthorURL:   c.AuthorURL,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Content:     c.Content,
		Created:     time.Now().Unix(),
		State:       models.CommentStateWaiting,
		IP:          r.RemoteAddr,
		Agent:       r.UserAgent(),
		IsAdmin:     false,
	}
	if _, err := db.Insert(comm); err != nil {
		logs.Error("postComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, nil, nil)
}

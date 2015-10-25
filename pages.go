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
	"github.com/caixw/typing/themes"
	"github.com/issue9/conv"
	"github.com/issue9/is"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

func getTagPosts(page int, tagID int64) ([]*themes.Post, error) {
	posts := make([]*themes.Post, 0, opt.PageSize)
	sql := `SELECT p.{id} AS ID, p.{name} AS Name,
		p.{title} AS Title, p.{summary} AS Summary, p.{created} AS Created, p.{allowComment} AS AllowComment
		FROM #relationships AS r
		LEFT JOIN #posts AS p ON p.{id}=r.{postID}
		WHERE p.{state}=? AND r.{tagID}=?
		ORDER BY {order} DESC
		LIMIT ? OFFSET ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, tagID, opt.PageSize, opt.PageSize*page)
	if err != nil {
		return nil, err
	}
	_, err = fetch.Obj(&posts, rows)
	rows.Close()

	return posts, err
}

func getPosts(page int) ([]*themes.Post, error) {
	posts := make([]*themes.Post, 0, opt.PageSize)
	sql := `SELECT {id} AS ID, {name} AS Name, {title} AS Title, {summary} AS Summary, {created} AS Created, {allowComment} AS AllowComment
	FROM #posts
	WHERE {state}=?
	ORDER BY {order} DESC
	LIMIT ? OFFSET ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, opt.PageSize, opt.PageSize*page)
	if err != nil {
		return nil, err
	}
	_, err = fetch.Obj(&posts, rows)
	rows.Close()

	return posts, err
}

// 首页或是列表页
func pagePosts(w http.ResponseWriter, r *http.Request) {
	info, err := themes.GetInfo()
	if err != nil {
		logs.Error("pagePosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	page := conv.MustInt(r.FormValue("page"), 1)
	if page < 1 { // 不能小于1
		page = 1
	}
	posts, err := getPosts(page - 1)
	if err != nil {
		logs.Error("pagePosts:", err)
		// TODO 显示一个正常的500页面，而不是json格式的
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	data := map[string]interface{}{
		"info":  info,
		"posts": posts,
	}
	themes.Render(w, "list", data)
}

func pageTags(w http.ResponseWriter, r *http.Request) {
	info, err := themes.GetInfo()
	if err != nil {
		logs.Error("pageTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	info.Canonical = opt.SiteURL + "tags"
	info.Title = "标签"

	themes.Render(w, "tags", info)
}

func pageTag(w http.ResponseWriter, r *http.Request) {
	info, err := themes.GetInfo()
	if err != nil {
		logs.Error("pageTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	tagName, ok := core.ParamString(w, r, "id")
	if !ok {
		return
	}
	tag := &models.Tag{Name: tagName}
	if err := db.Select(tag); err != nil {
		logs.Error("pageTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	page := conv.MustInt(r.FormValue("page"), 1)
	if page < 1 { // 不能小于1
		page = 1
	}
	posts, err := getTagPosts(page-1, tag.ID)
	if err != nil {
		logs.Error("pageTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	data := map[string]interface{}{
		"info":  info,
		"tag":   tag,
		"posts": posts,
	}
	themes.Render(w, "tag", data)
}

func pagePost(w http.ResponseWriter, r *http.Request) {
	idStr, ok := core.ParamString(w, r, "id")
	idStr = strings.TrimSuffix(idStr, opt.Suffix)
	if !ok {
		return
	}

	var mp *models.Post
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		mp = &models.Post{Name: idStr}
	} else {
		mp = &models.Post{ID: postID}
	}
	if err := db.Select(mp); err != nil {
		logs.Error("pagePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if mp.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	// commentsSize
	rs := &models.Relationship{PostID: mp.ID}
	commentsSize, err := db.Count(rs)
	if err != nil {
		logs.Error("pagePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	post := &themes.Post{
		ID:           mp.ID,
		Name:         mp.Name,
		Title:        mp.Title,
		Summary:      mp.Summary,
		Content:      mp.Content,
		Author:       opt.ScreenName,
		CommentsSize: commentsSize,
		Created:      mp.Created,
		Modified:     mp.Modified,
		AllowComment: mp.AllowComment,
	}

	info, err := themes.GetInfo()
	if err != nil {
		logs.Error("pagePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	data := map[string]interface{}{
		"info": info,
		"post": post,
	}
	themes.Render(w, "post", data)
}

// @api get /api/posts/{id}/comments
// @apiQuery page  int 页码
// @apiGroup front
//
// @apiSuccess 200 OK
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

	var page int
	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	sql.Limit(opt.PageSize, page*opt.PageSize)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("getComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"comments": maps}, nil)
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

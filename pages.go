// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/themes"
	"github.com/issue9/conv"
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

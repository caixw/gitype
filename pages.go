// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/themes"
	"github.com/issue9/conv"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

func getPageInfo() (*themes.PageInfo, error) {
	p := &themes.PageInfo{
		SiteName:    opt.SiteName,
		SecondTitle: opt.SecondTitle,
		Keywords:    opt.Keywords,
		Description: opt.Description,
		AppVersion:  version,
		GoVersion:   runtime.Version(),
	}

	var err error
	sql := "SELECT COUNT(*) as cnt FROM #posts WHERE {state}=?" // TODO 预编译成stmt
	if p.PostSize, err = getSize(sql, models.PostStatePublished); err != nil {
		return nil, err
	}

	sql = "SELECT COUNT(*) as cnt FROM #comments WHERE {state}=?"
	if p.CommentSize, err = getSize(sql, models.CommentStateApproved); err != nil {
		return nil, err
	}

	if p.Tags, err = getTags(); err != nil {
		return nil, err
	}

	// Topics
	sql = `SELECT c.{content}, p.{title}, p.{name}, p.{id}
	FROM #comments AS c
	LEFT JOIN #posts AS p ON c.{postID}=p.{id}
	ORDER BY c.{id} DESC
	LIMIT ? `
	rows, err := db.Query(true, sql, opt.SidebarSize)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	p.Topics = make([]themes.Anchor, 0, opt.SidebarSize)
	for _, v := range maps {
		a := themes.Anchor{
			Title: v["content"],
			Ext:   v["title"],
		}
		if len(v["name"]) > 0 {
			a.Link = "/posts/" + v["name"] + opt.Suffix
		} else {
			a.Link = "/posts/" + v["id"] + opt.Suffix
		}
		p.Topics = append(p.Topics, a)
	}

	return p, nil
}

func getTags() ([]themes.Anchor, error) {
	sql := "SELECT {id}, {title}, {name}, {description} FROM #tags"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	ret := make([]themes.Anchor, 0, len(maps))
	for _, v := range maps {
		a := themes.Anchor{
			Title: v["title"],
			Ext:   v["description"],
			Link:  core.TagURL(opt, v["name"]),
		}
		ret = append(ret, a)
	}

	return ret, nil
}

func getSize(sql string, args ...interface{}) (int, error) {
	rows, err := db.Query(true, sql, args...)
	if err != nil {
		return 0, err
	}
	cnts, err := fetch.ColumnString(true, "cnt", rows)
	rows.Close()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cnts[0])
}

func getPosts(page int) ([]*themes.Post, error) {
	posts := make([]*themes.Post, 0, opt.PageSize)
	sql := `SELECT {title} AS Title, {summary} AS Summary, {created} AS Created, {allowComment} AS AllowComment
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
	p, err := getPageInfo()
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
		"page":  p,
		"posts": posts,
	}
	themes.Render(w, "list", data)
}

func pageTags(w http.ResponseWriter, r *http.Request) {
}

func pageTag(w http.ResponseWriter, r *http.Request) {

}

func pagePost(w http.ResponseWriter, r *http.Request) {

}

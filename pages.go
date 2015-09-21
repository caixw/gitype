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
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// 页面的基本信息
type page struct {
	Title       string
	SiteName    string
	SecondTitle string
	Keywords    string
	Description string
	AppVersion  string
	GoVersion   string

	PostSize    int      // 文章数量
	CommentSize int      // 评论数量
	Tags        []anchor // 标签列表
	Cats        []anchor // 分类列表
	Topics      []anchor // 最新评论的10条内容
}

type anchor struct {
	Link  string // 链接地址
	Title string // 地址的字面文字
	Ext   string // 扩展内容，比如title,alt等，根据链接来确定
}

func getPage() (*page, error) {
	p := &page{
		SiteName:    opt.SiteName,
		SecondTitle: opt.SecondTitle,
		Keywords:    opt.Keywords,
		Description: opt.Description,
		AppVersion:  version,
		GoVersion:   runtime.Version(),
	}

	var err error
	sql := "SELECT COUNT(*) as cnt FROM #posts WHERE {state}=? AND !password"
	if p.PostSize, err = getSize(sql, models.PostStateNormal); err != nil {
		return nil, err
	}

	sql = "SELECT COUNT(*) as cnt FROM #comments WHERE {state}=?"
	if p.CommentSize, err = getSize(sql, models.CommentStateApproved); err != nil {
		return nil, err
	}

	if p.Tags, err = getMetas(models.MetaTypeTag); err != nil {
		return nil, err
	}

	if p.Cats, err = getMetas(models.MetaTypeCat); err != nil {
		return nil, err
	}

	// Topics
	sql = `SELECT c.{content}, p.{title}, p.{name}, p.{id}
	FROM #comments AS c
	LEFT JOIN #posts AS p ON c.{postID}=p.{id}
	LIMIT ? ORDER BY c.{id} DESC`
	rows, err := db.Query(true, sql, opt.SidebarSize)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	p.Topics = make([]anchor, 0, opt.SidebarSize)
	for _, v := range maps {
		a := anchor{
			Title: v["content"],
			Ext:   v["title"],
		}
		if len(v["name"]) > 0 {
			a.Link = "/posts/" + v["name"]
		} else {
			a.Link = "/posts/" + v["id"]
		}
		p.Topics = append(p.Topics, a)
	}

	return p, nil
}

func getMetas(typ int) ([]anchor, error) {
	sql := "SELECT {id}, {title}, {name}, {description} FROM #metas WHERE {type}=?"
	rows, err := db.Query(true, sql, typ)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	var link string
	if typ == models.MetaTypeCat {
		link = "/cats/"
	} else {
		link = "/tags/"
	}

	ret := make([]anchor, 0, len(maps))
	for _, v := range maps {
		a := anchor{
			Title: v["title"],
			Ext:   v["description"],
		}
		if len(v["name"]) > 0 {
			a.Link = link + v["name"]
		} else {
			a.Link = link + v["id"]
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

func pageIndex(w http.ResponseWriter, r *http.Request) {
	p, err := getPage()
	if err != nil {
		logs.Error("pageIndex:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := map[string]interface{}{
		"page": p,
	}
	if err := themes.Render(w, "index", data); err != nil {
		logs.Error("pageIndex:", err)
	}
}

func pageTags(w http.ResponseWriter, r *http.Request) {

}

func pageTag(w http.ResponseWriter, r *http.Request) {

}

func pageCats(w http.ResponseWriter, r *http.Request) {

}

func pageCat(w http.ResponseWriter, r *http.Request) {

}

func pagePosts(w http.ResponseWriter, r *http.Request) {

}

func pagePost(w http.ResponseWriter, r *http.Request) {

}

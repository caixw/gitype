// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/issue9/logs"
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

func pageIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"page": &page{},
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

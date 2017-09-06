// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"strings"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
)

// /search.html?q=key&page=2
func (client *Client) getSearch(w http.ResponseWriter, r *http.Request) {
	p := client.page(typeSearch)

	q := r.FormValue("q")
	if len(q) == 0 {
		http.Redirect(w, r, vars.PostsURL(1), http.StatusPermanentRedirect)
		return
	}

	page, ok := client.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("参数 page: %d 小于 1", page)
		client.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	// 获取所有的搜索结果
	posts := search(q, client.data)

	p.Title = "搜索:" + q
	p.Q = q
	p.Keywords = q + ",搜索,search"
	p.Description = "搜索关键字" + q + "的结果"
	p.Canonical = client.data.URL(vars.SearchURL(p.Q, page))
	start, end, ok := client.getPostsRange(len(posts), page, w)
	if !ok {
		return
	}
	p.Posts = posts[start:end]
	if page > 1 {
		p.prevPage(vars.SearchURL(q, page-1), "")
	}
	if end < len(posts) {
		p.nextPage(vars.SearchURL(q, page+1), "")
	}

	p.render(w, "search", nil)
}

// 查找出所有符合要求的文章列表
func search(q string, d *data.Data) []*data.Post {
	index := strings.IndexByte(q, ':')
	if index <= 0 {
		return searchDefault(q, d)
	}

	index++
	typ := q[:index]
	content := strings.TrimSpace(q[index:])

	switch typ {
	case "date:":
		return searchDate(content, d)
	case "tag:":
		return searchTag(content, d)
	case "title:":
		return searchTitle(content, d)
	}

	// 不存在的分类，则使用全部文字按默认情况进行搜索
	return searchDefault(q, d)
}

// 按日期进行分类
func searchDate(date string, d *data.Data) []*data.Post {
	// TODO
}

// 按标签进行搜索
func searchTag(q string, d *data.Data) []*data.Post {
	// TODO
}

// 仅搜索标题
func searchTitle(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))
	key := strings.ToLower(q)

	for _, v := range d.Posts {
		if strings.Contains(v.Title, key) {
			posts = append(posts, v)
		}
	}

	return posts
}

// 默认情况下，搜索标题和内容
func searchDefault(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))
	key := strings.ToLower(q)

	for _, v := range d.Posts {
		if strings.Contains(v.Title, key) || strings.Contains(v.Content, key) {
			posts = append(posts, v)
		}
	}

	return posts
}

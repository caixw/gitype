// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/vars"
	"github.com/caixw/gitype/vars/url"
	"github.com/issue9/logs"
)

// /search.html?q=key&page=2
func (client *Client) getSearch(w http.ResponseWriter, r *http.Request) {
	p := client.page(vars.PageSearch, w, r)

	q := r.FormValue(vars.URLQueryQ)
	if len(q) == 0 {
		http.Redirect(w, r, url.Posts(1), http.StatusPermanentRedirect)
		return
	}

	page, ok := client.queryInt(w, r, vars.URLQueryPage, 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("参数 page: %d 小于 1", page)
		client.renderError(w, r, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	pp := client.data.Pages[vars.PageSearch]
	p.Title = helper.ReplaceContent(pp.Title, q)
	p.Keywords = helper.ReplaceContent(pp.Keywords, q)
	p.Description = helper.ReplaceContent(pp.Description, q)
	p.Q = q
	p.Canonical = client.data.BuildURL(url.Search(p.Q, page))

	posts := search(q, client.data) // 获取所有的搜索结果
	start, end, ok := client.getPostsRange(len(posts), page, w, r)
	if !ok {
		return
	}
	p.Posts = posts[start:end]
	if page > 1 {
		p.prevPage(url.Search(q, page-1), "")
	}
	if end < len(posts) {
		p.nextPage(url.Search(q, page+1), "")
	}

	p.render(vars.PageSearch)
}

// 查找出所有符合要求的文章列表
func search(q string, d *data.Data) []*data.Post {
	index := strings.IndexByte(q, vars.SearchKeySeparator)
	// 若 : 前后为空，则直接将整个字符串当作搜索关键字
	if index <= 0 || len(q)-1 == index {
		return searchDefault(q, d)
	}

	typ := q[:index]
	content := strings.ToLower(strings.TrimSpace(q[index+1:]))

	switch typ {
	case vars.SearchKeyTag:
		return searchTag(content, d)
	case vars.SearchKeySeries:
		return searchSeries(content, d)
	case vars.SearchKeyTitle:
		return searchTitle(content, d)
	}

	// 不存在的分类，则使用全部文字按默认情况进行搜索
	return searchDefault(q, d)
}

// 按标签进行搜索
func searchSeries(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))

	for _, tag := range d.Series {
		if strings.Contains(strings.ToLower(tag.Title), q) {
			posts = append(posts, tag.Posts...)
		}
	}

	return posts
}

// 按标签进行搜索
func searchTag(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))

	for _, tag := range d.Tags {
		if strings.Contains(strings.ToLower(tag.Title), q) {
			posts = append(posts, tag.Posts...)
		}
	}

	return posts
}

// 仅搜索标题
func searchTitle(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))

	for _, post := range d.Posts {
		if strings.Contains(strings.ToLower(post.Title), q) {
			posts = append(posts, post)
		}
	}

	fmt.Println(len(posts))
	return posts
}

// 默认情况下，搜索标题和内容
func searchDefault(q string, d *data.Data) []*data.Post {
	posts := make([]*data.Post, 0, len(d.Posts))

	for _, post := range d.Posts {
		if strings.Contains(post.Title, q) || strings.Contains(post.Content, q) {
			posts = append(posts, post)
		}
	}

	return posts
}

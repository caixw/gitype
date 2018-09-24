// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"

	"github.com/issue9/logs"
	"github.com/issue9/web"
	"github.com/issue9/web/context"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/vars"
)

func (client *Client) initRoutes() error {
	err := client.initFeedRoutes()
	if err != nil {
		return err
	}

	handle := func(pattern string, h http.HandlerFunc) {
		if err != nil {
			return
		}

		client.patterns = append(client.patterns, pattern)
		err = client.mux.HandleFunc(pattern, client.prepare(h), http.MethodGet)
	}

	handle(vars.PostURL("{slug}"), client.getPost)                 // posts/2016/about.html   posts/{slug}.html
	handle(vars.AssetURL("{path}"), client.getAsset)               // posts/2016/about/abc.png  posts/{path}
	handle(vars.IndexURL(0), client.getPosts)                      // index.html
	handle(vars.LinksURL(), client.getLinks)                       // links.html
	handle(vars.TagURL("{slug}", 1), client.getTag)                // tags/tag1.html     tags/{slug}.html
	handle(vars.TagsURL(), client.getTags)                         // tags.html
	handle(vars.ArchivesURL(), client.getArchives)                 // archives.html
	handle(vars.SearchURL("", 1), client.getSearch)                // search.html
	handle(vars.ThemeURL("{path}"), client.getTheme)               // themes/...          themes/{path}
	handle("/{path}", client.getRaw)                               // /...                /{path}
	handle(client.data.ServiceWorkerPath, client.getServiceWorker) // /sw.js

	return err
}

func (client *Client) initFeedRoutes() (err error) {
	handle := func(feed *data.Feed) {
		if err != nil || feed == nil {
			return
		}

		client.patterns = append(client.patterns, feed.URL)
		err = client.mux.HandleFunc(feed.URL, client.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", feed.Type)
			w.Write(feed.Content)
		}), http.MethodGet)
	}

	handle(client.data.RSS)
	handle(client.data.Atom)
	handle(client.data.Sitemap)
	handle(client.data.Opensearch)
	handle(client.data.Manifest)

	return err
}

// /sw.js
func (client *Client) getServiceWorker(w http.ResponseWriter, r *http.Request) {
	// https://github.com/golang/go/issues/17083
	// 需要保证 Header().Set 在 WriteHeader 之前调用
	w.Header().Set("Content-Type", "application/javascript")
	w.WriteHeader(http.StatusOK)
	w.Write(client.data.ServiceWorker)
}

// 文章详细页
// /posts/{slug}.html
func (client *Client) getPost(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	slug, err := ctx.ParamString("slug")
	if err != nil {
		logs.Error(err)
		client.getAsset(w, r)
		return
	}

	index := -1
	for i, p := range client.data.Posts {
		if p.Slug == slug {
			index = i
			break
		}
	}

	if index < 0 {
		logs.Debugf("并未找到与之相对应的文章：%s", slug)
		client.getRaw(w, r) // 文章不存在，则查找 raws 目录下是否存在同名文件
		return
	}

	post := client.data.Posts[index]
	p := client.page(ctx, vars.PagePost)

	p.Post = post
	p.Keywords = post.Keywords
	p.Description = post.Summary
	p.Title = post.HTMLTitle
	p.Canonical = web.URL(post.Permalink)
	p.License = post.License // 文章可具体指定协议
	p.Author = post.Author   // 文章可具体指定作者

	if index > 0 {
		prev := client.data.Posts[index-1]
		p.Prev(prev.Permalink, prev.Title)
	}
	if index+1 < len(client.data.Posts) {
		next := client.data.Posts[index+1]
		p.Next(next.Permalink, next.Title)
	}

	p.Render(post.Template)
}

// 首页及文章列表页
// /
// /index.html?page=2
func (client *Client) getPosts(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	page := client.queryInt(ctx, vars.URLQueryPage, 1)

	if page < 1 {
		ctx.Exit(http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := client.page(ctx, vars.PageIndex)
	if page > 1 { // 非首页，标题显示页码数
		p.Type = vars.PagePosts
	}
	pp := client.data.Pages[vars.PagePosts]
	p.Title = pp.Title
	p.Keywords = pp.Keywords
	p.Description = pp.Description
	p.Canonical = web.URL(vars.PostsURL(page))

	start, end, ok := client.getPostsRange(len(client.data.Posts), page, w, r)
	if !ok {
		return
	}
	p.Posts = client.data.Posts[start:end]
	if page > 1 {
		p.Prev(vars.PostsURL(page-1), "")
	}
	if end < len(client.data.Posts) {
		p.Next(vars.PostsURL(page+1), "")
	}

	p.Render(vars.PagePosts)
}

// 标签详细页
// /tags/tag1.html?page=2
func (client *Client) getTag(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	slug, err := ctx.ParamString("slug")
	if err != nil {
		logs.Error(err)
		client.getRaw(w, r)
		return
	}

	var tag *data.Tag
	for _, t := range client.data.Tags {
		if t.Slug == slug {
			tag = t
			break
		}
	}

	if tag == nil {
		logs.Debugf("查找的标签 %s 不存在", slug)
		client.getRaw(w, r) // 标签不存在，则查找该文件是否存在于 raws 目录下。
		return
	}

	page := client.queryInt(ctx, vars.URLQueryPage, 1)
	if page < 1 {
		ctx.Exit(http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
	}

	p := client.page(ctx, vars.PageTag)
	p.Tag = tag
	p.Title = tag.HTMLTitle
	p.Keywords = tag.Keywords
	p.Description = tag.Content
	p.Canonical = web.URL(vars.TagURL(slug, page))

	start, end, ok := client.getPostsRange(len(tag.Posts), page, w, r)
	if !ok {
		return
	}
	p.Posts = tag.Posts[start:end]
	if page > 1 {
		p.Prev(vars.TagURL(slug, page-1), "")
	}
	if end < len(tag.Posts) {
		p.Next(vars.TagURL(slug, page+1), "")
	}

	p.Render(vars.PageTag)
}

// 友情链接页
// /links.html
func (client *Client) getLinks(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	p := client.page(ctx, vars.PageLinks)
	pp := client.data.Pages[vars.PageLinks]
	p.Title = pp.Title
	p.Keywords = pp.Keywords
	p.Description = pp.Description
	p.Canonical = web.URL(vars.LinksURL())

	p.Render(vars.PageLinks)
}

// 标签列表页
// /tags.html
func (client *Client) getTags(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	p := client.page(ctx, vars.PageTags)
	pp := client.data.Pages[vars.PageTags]
	p.Title = pp.Title
	p.Keywords = pp.Keywords
	p.Description = pp.Description
	p.Canonical = web.URL(vars.TagsURL())

	p.Render(vars.PageTags)
}

// 归档页
// /archives.html
func (client *Client) getArchives(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	p := client.page(ctx, vars.PageArchives)
	pp := client.data.Pages[vars.PageArchives]
	p.Title = pp.Title
	p.Keywords = pp.Keywords
	p.Description = pp.Description
	p.Canonical = web.URL(vars.ArchivesURL())
	p.Archives = client.data.Archives

	p.Render(vars.PageArchives)
}

// 确认当前文章列表页选择范围。
func (client *Client) getPostsRange(postsSize, page int, w http.ResponseWriter, r *http.Request) (start, end int, ok bool) {
	ctx := web.NewContext(w, r)
	size := client.data.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%d]，实际文章数量为[%d]\n", page, postsSize)
		ctx.Exit(http.StatusNotFound) // 页码超出范围，不存在
		return 0, 0, false
	}

	end = start + size
	if postsSize < end {
		end = postsSize
	}

	return start, end, true
}

// 获取查询参数 key 的值，并将其转换成 Int 类型，若该值不存在返回 def 作为其默认值，
// 若是类型不正确，则返回一个 false，并向客户端输出一个 400 错误。
func (client *Client) queryInt(ctx *context.Context, key string, def int) int {
	q := ctx.Queries()
	v := q.Int(key, def)

	if q.HasErrors() {
		logs.Error(q.Errors()[key])
		ctx.Exit(http.StatusBadRequest)
	}

	return v
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/handlers"
	"github.com/issue9/logs"
)

// /
func (a *app) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == a.data.URLS.Root || r.URL.Path == a.data.URLS.Root+"/" {
		a.getPosts(w, r)
		return
	}

	root := http.Dir(a.path.DataRaws)
	prefix := a.data.URLS.Root + "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// /search.html?q=key&page=2
func (a *app) getSearch(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()

	key := r.FormValue("q")
	if len(key) == 0 {
		p.render(w, r, "search", map[string]string{"Content-Type": "text/html"})
		return
	}

	page, ok := queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1\n", page)
		w.WriteHeader(http.StatusNotFound) // 页码为负数的表示不存在
		return
	}

	// 查找标题和内容
	posts := make([]*data.Post, 0, 10)
	for _, v := range a.data.Posts {
		if strings.Index(v.Title, key) >= 0 || strings.Index(v.Content, key) >= 0 {
			posts = append(posts, v)
		}
	}

	p.Title = "搜索:" + key
	p.Keywords = key
	p.Description = "搜索关键字" + key + "的结果集"
	start, end, ok := a.getPostsRange(len(posts), page, w)
	if !ok {
		return
	}
	p.Posts = posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  a.searchURL(key, page-1), // 页码从1开始计数
		}
	}
	if end < len(posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  a.searchURL(key, page+1),
		}
	}

	p.render(w, r, "search", map[string]string{"Content-Type": "text/html"})
}

// 确认当前文章列表页选择范围。
func (a *app) getPostsRange(postsSize, page int, w http.ResponseWriter) (start, end int, ok bool) {
	size := a.data.Config.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%v]，实际文章数量为[%v]\n", page, postsSize)
		w.WriteHeader(http.StatusNotFound) // 页码超出范围，不存在
		return 0, 0, false
	}

	end = start + size
	if postsSize < end {
		end = postsSize
	}

	return start, end, true
}

// 获取媒体文件
//
// /media/2015/intro-php/content.html ==> /posts/2015/intro-php/content.html
func (a *app) getMedia(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix(vars.MediaURL, http.FileServer(http.Dir(a.path.DataPosts))).ServeHTTP(w, r)
}

// 首页及文章列表页
// /
// /posts.html?page=2
func (a *app) getPosts(w http.ResponseWriter, r *http.Request) {
	page, ok := queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1\n", page)
		w.WriteHeader(http.StatusNotFound) // 页码为负数的表示不存在
		return
	}

	p := a.newPage()
	if page > 1 { // 非首页，标题显示页码数
		p.Title = fmt.Sprintf("第%v页", page)
	}
	p.Canonical = a.postsURL(page)
	start, end, ok := a.getPostsRange(len(a.data.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = a.data.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  a.postsURL(page - 1), // 页码从1开始计数
		}
	}
	if end < len(a.data.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  a.postsURL(page + 1),
		}
	}

	p.render(w, r, "posts", map[string]string{"Content-Type": "text/html"})
}

// 主题文件
// /themes/...
func (a *app) getThemes(w http.ResponseWriter, r *http.Request) {
	root := http.Dir(a.path.DataThemes)
	prefix := a.data.URLS.Root + a.data.URLS.Themes
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// 标签详细页
// /tags/tag1.html?page=2
func (a *app) getTag(w http.ResponseWriter, r *http.Request) {
	slug, ok := paramString(w, r, "slug")
	if !ok {
		return
	}

	var tag *data.Tag
	for _, t := range a.data.Tags {
		if t.Slug == slug {
			tag = t
			break
		}
	}
	if tag == nil {
		logs.Debugf("查找的标签[%v]不存在\n", slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page, ok := queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1\n", page)
		w.WriteHeader(http.StatusNotFound) // 页码为负数的页码不存在
		return
	}

	p := a.newPage()
	p.Tag = tag
	p.Title = tag.Title
	p.Canonical = a.tagURL(slug, page)
	p.Description = "标签" + tag.Title + "的介绍"

	start, end, ok := a.getPostsRange(len(tag.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = tag.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  a.tagURL(slug, page-1), // 页码从1开始计数
		}
	}
	if end < len(tag.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  a.tagURL(slug, page+1), // 页码从1开始计数
		}
	}
	p.render(w, r, "tag", map[string]string{"Content-Type": "text/html"})
}

// 标签列表页
// /tags.html
func (a *app) getTags(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()
	p.Title = "标签"
	p.Canonical = a.tagsURL()

	p.render(w, r, "tags", map[string]string{"Content-Type": "text/html"})
}

// 文章详细页
// /posts/{slug}.html
func (a *app) getPost(w http.ResponseWriter, r *http.Request) {
	slug, ok := paramString(w, r, "slug")
	if !ok {
		return
	}

	var post *data.Post
	var next, prev *data.Link
	for index, p := range a.data.Posts {
		if p.Slug != slug {
			continue
		}
		post = p

		if index > 0 {
			p := a.data.Posts[index-1]
			prev = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}

		index++
		if index < len(a.data.Posts) {
			p := a.data.Posts[index]
			next = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}
	} // end for a.data.Posts

	if post == nil {
		logs.Debugf("并未找到与之相对应的文章\n", slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p := a.newPage()
	p.Post = post
	p.NextPage = next
	p.PrevPage = prev
	p.Canonical = post.Permalink
	p.Description = post.Summary
	p.Title = post.Title
	if len(p.Tags) > 0 {
		keywords := make([]string, 0, len(p.Tags))
		for _, v := range p.Tags {
			keywords = append(keywords, v.Title)
		}
		p.Keywords = strings.Join(keywords, ",")
	}

	p.render(w, r, post.Template, map[string]string{"Content-Type": "text/html"})
}

// 每次访问前需要做的预处理工作。
func (a *app) pre(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.conf.isDebug() { // 调试状态，则每次都重新加载数据
			if err := a.reload(); err != nil {
				logs.Error("app.pre:", err)
			}
		}

		// 输出访问日志
		logs.Infof("%v：%v\n", r.UserAgent(), r.URL)

		// 直接根据整个博客的最后更新时间来确认etag
		if r.Header.Get("If-None-Match") == a.etag {
			logs.Infof("304:%v\n", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", a.etag)
		handlers.CompressFunc(f).ServeHTTP(w, r)
	}
}

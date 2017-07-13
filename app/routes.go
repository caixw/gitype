// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/middleware/compress"
	"github.com/issue9/utils"
)

func (a *app) initRoutes() error {
	handle := func(pattern string, h http.HandlerFunc) error {
		return a.mux.HandleFunc(pattern, a.prepare(h), http.MethodGet)
	}

	// posts/2016/about.html
	if err := handle(vars.PostURL("{slug}"), a.getPost); err != nil {
		return err
	}

	// index.html
	pattern := vars.Posts + vars.Suffix
	if err := handle(pattern, a.getPosts); err != nil {
		return err
	}

	// links.html
	pattern = vars.Links + vars.Suffix
	if err := handle(pattern, a.getLinks); err != nil {
		return err
	}

	// tags/tag1.html
	if err := handle(vars.TagURL("{slug}", 1), a.getTag); err != nil {
		return err
	}

	// tags.html
	pattern = vars.Tags + vars.Suffix
	if err := handle(pattern, a.getTags); err != nil {
		return err
	}

	// search.html
	if err := handle(vars.SearchURL("", 1), a.getSearch); err != nil {
		return err
	}

	// themes/...
	pattern = vars.Themes + "/{path}"
	if err := handle(pattern, a.getThemes); err != nil {
		return err
	}

	// /...
	pattern = "/{path}"
	return handle(pattern, a.getRaws)
}

// 文章详细页
// /posts/{slug}.html
func (a *app) getPost(w http.ResponseWriter, r *http.Request) {
	id, found := a.paramString(w, r, "slug")
	if !found {
		return
	}

	var post *data.Post
	var next, prev *data.Link
	for index, p := range a.buf.Data.Posts {
		if p.Slug != id {
			continue
		}
		post = p

		if index > 0 {
			p := a.buf.Data.Posts[index-1]
			prev = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}

		index++
		if index < len(a.buf.Data.Posts) {
			p := a.buf.Data.Posts[index]
			next = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}
	} // end for a.data.Posts

	if post == nil {
		logs.Debugf("并未找到与之相对应的文章:%v", id)
		a.getRaws(w, r) // 文章不存在，则查找 raws 目录下是否存在同名文件
		return
	}

	p := a.newPage()
	p.Post = post
	p.NextPage = next
	p.PrevPage = prev
	p.Keywords = post.Keywords
	p.Description = post.Summary
	p.Title = post.Title
	p.Canonical = post.Permalink

	p.render(w, post.Template, nil)
}

// 首页及文章列表页
// /
// /posts.html?page=2
func (a *app) getPosts(w http.ResponseWriter, r *http.Request) {
	page, ok := a.queryInt(w, r, "page", 1)
	if !ok {
		return
	}

	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1\n", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := a.newPage()
	if page > 1 { // 非首页，标题显示页码数
		p.Title = fmt.Sprintf("第%v页", page)
	}
	p.Canonical = vars.PostsURL(page)

	start, end, ok := a.getPostsRange(len(a.buf.Data.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = a.buf.Data.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  vars.PostsURL(page - 1), // 页码从 1 开始计数
		}
	}
	if end < len(a.buf.Data.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  vars.PostsURL(page + 1),
		}
	}

	p.render(w, "posts", nil)
}

// 标签详细页
// /tags/tag1.html?page=2
func (a *app) getTag(w http.ResponseWriter, r *http.Request) {
	slug, ok := a.paramString(w, r, "slug")
	if !ok {
		return
	}

	var tag *data.Tag
	for _, t := range a.buf.Data.Tags {
		if t.Slug == slug {
			tag = t
			break
		}
	}

	if tag == nil {
		logs.Debugf("查找的标签[%v]不存在", slug)
		a.getRaws(w, r) // 标签不存在，则查找该文件是否存在于 raws 目录下。
		return
	}

	page, ok := a.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := a.newPage()
	p.Tag = tag
	p.Title = tag.Title
	p.Canonical = vars.TagURL(slug, page)
	p.Description = "标签" + tag.Title + "的介绍"

	start, end, ok := a.getPostsRange(len(tag.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = tag.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  vars.TagURL(slug, page-1), // 页码从1开始计数
		}
	}
	if end < len(tag.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  vars.TagURL(slug, page+1), // 页码从1开始计数
		}
	}

	p.render(w, "tag", nil)
}

// 友情链接页
// /links.html
func (a *app) getLinks(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()
	p.Title = "友情链接"
	p.Canonical = vars.Links + vars.Suffix

	p.render(w, "links", nil)
}

// 标签列表页
// /tags.html
func (a *app) getTags(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()
	p.Title = "标签"
	p.Canonical = vars.Tags + vars.Suffix

	p.render(w, "tags", nil)
}

// 主题文件
// /themes/...
func (a *app) getThemes(w http.ResponseWriter, r *http.Request) {
	root := http.Dir(a.path.ThemesDir)
	http.StripPrefix(vars.Themes, http.FileServer(root)).ServeHTTP(w, r)
}

// /search.html?q=key&page=2
func (a *app) getSearch(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()

	key := r.FormValue("q")
	if len(key) == 0 {
		p.render(w, "search", nil)
		return
	}

	page, ok := a.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("参数 page: %v 小于 1", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	// 查找标题和内容
	posts := make([]*data.Post, 0, a.buf.Data.Config.PageSize)
	for _, v := range a.buf.Data.Posts {
		if strings.Index(v.Title, key) >= 0 || strings.Index(v.Content, key) >= 0 {
			posts = append(posts, v)
		}
	}

	p.Title = "搜索:" + key
	p.Q = key
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
			URL:  vars.SearchURL(key, page-1), // 页码从1开始计数
		}
	}
	if end < len(posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  vars.SearchURL(key, page+1),
		}
	}

	p.render(w, "search", nil)

}

// 读取根下的文件
// /...
func (a *app) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		a.getPosts(w, r)
		return
	}

	root := http.Dir(a.path.RawsDir)
	if !utils.FileExists(filepath.Join(a.path.RawsDir, r.URL.Path)) {
		a.renderError(w, http.StatusNotFound)
		return
	}
	prefix := "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)

}

// 确认当前文章列表页选择范围。
func (a *app) getPostsRange(postsSize, page int, w http.ResponseWriter) (start, end int, ok bool) {
	size := a.buf.Data.Config.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%v]，实际文章数量为[%v]\n", page, postsSize)
		a.renderError(w, http.StatusNotFound) // 页码超出范围，不存在
		return 0, 0, false
	}

	end = start + size
	if postsSize < end {
		end = postsSize
	}

	return start, end, true
}

// 每次访问前需要做的预处理工作。
func (a *app) prepare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs.Infof("%v: %v", r.UserAgent(), r.URL) // 输出访问日志

		// 直接根据整个博客的最后更新时间来确认 etag
		if r.Header.Get("If-None-Match") == a.buf.Etag {
			logs.Infof("304: %v", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", a.buf.Etag)
		compress.New(f, logs.ERROR()).ServeHTTP(w, r)
	}
}

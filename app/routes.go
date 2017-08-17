// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/middleware/compress"
	"github.com/issue9/utils"
)

// 模板的扩展名，在主题目录下，所有该扩展名的文件，不会被展示
var (
	ignoreThemeFileExts = []string{
		".html",
		".yaml",
		".yml",
	}
)

func isIgnoreThemeFile(file string) bool {
	ext := filepath.Ext(file)

	for _, v := range ignoreThemeFileExts {
		if ext == v {
			return true
		}
	}

	return false
}

func (a *app) initRoutes() error {
	var err error
	handle := func(pattern string, h http.HandlerFunc) {
		if err != nil {
			return
		}
		err = a.mux.HandleFunc(pattern, a.prepare(h), http.MethodGet)
	}

	handle(vars.PostURL("{slug}"), a.getPost)     // posts/2016/about.html   posts/{slug}.html
	handle(vars.IndexURL(0), a.getPosts)          // index.html
	handle(vars.LinksURL(), a.getLinks)           // links.html
	handle(vars.TagURL("{slug}", 1), a.getTag)    // tags/tag1.html     tags/{slug}.html
	handle(vars.TagsURL(), a.getTags)             // tags.html
	handle(vars.SearchURL("", 1), a.getSearch)    // search.html
	handle(vars.ThemesURL("{path}"), a.getThemes) // themes/...          themes/{path}
	handle("/{path}", a.getRaws)                  // /...                /{path}

	return err
}

// 文章详细页
// /posts/{slug}.html
func (a *app) getPost(w http.ResponseWriter, r *http.Request) {
	id, found := a.paramString(w, r, "slug")
	if !found {
		return
	}

	var index int
	for i, p := range a.buf.Data.Posts {
		if p.Slug == id {
			index = i
			break
		}
	}

	if index < 0 {
		logs.Debugf("并未找到与之相对应的文章:%s", id)
		a.getRaws(w, r) // 文章不存在，则查找 raws 目录下是否存在同名文件
		return
	}

	post := a.buf.Data.Posts[index]

	p := a.page(typePost)
	p.Post = post
	p.Keywords = post.Keywords
	p.Description = post.Summary
	p.Title = post.Title
	p.Canonical = post.Permalink

	if index > 0 {
		prev := a.buf.Data.Posts[index-1]
		p.prevPage(prev.Permalink, prev.Title)
	}
	if index+1 < len(a.buf.Data.Posts) {
		next := a.buf.Data.Posts[index+1]
		p.nextPage(next.Permalink, next.Title)
	}

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
		logs.Debugf("请求的页码[%d]小于1\n", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := a.page(typeIndex)
	if page > 1 { // 非首页，标题显示页码数
		p.Type = typePosts
		p.Title = fmt.Sprintf("第 %d 页", page)
	}
	p.Canonical = vars.PostsURL(page)

	start, end, ok := a.getPostsRange(len(a.buf.Data.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = a.buf.Data.Posts[start:end]
	if page > 1 {
		p.prevPage(vars.PostsURL(page-1), "")
	}
	if end < len(a.buf.Data.Posts) {
		p.nextPage(vars.PostsURL(page+1), "")
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
		logs.Debugf("查找的标签[%s]不存在", slug)
		a.getRaws(w, r) // 标签不存在，则查找该文件是否存在于 raws 目录下。
		return
	}

	page, ok := a.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%d]小于1", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := a.page(typeTag)
	p.Tag = tag
	p.Title = tag.Title
	p.Keywords = tag.Keywords
	p.Description = tag.Description
	p.Canonical = vars.TagURL(slug, page)

	start, end, ok := a.getPostsRange(len(tag.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = tag.Posts[start:end]
	if page > 1 {
		p.prevPage(vars.TagURL(slug, page-1), "")
	}
	if end < len(tag.Posts) {
		p.nextPage(vars.TagURL(slug, page+1), "")
	}

	p.render(w, "tag", nil)
}

// 友情链接页
// /links.html
func (a *app) getLinks(w http.ResponseWriter, r *http.Request) {
	p := a.page(typeLinks)
	p.Title = "友情链接"
	p.Canonical = vars.LinksURL()

	p.render(w, "links", nil)
}

// 标签列表页
// /tags.html
func (a *app) getTags(w http.ResponseWriter, r *http.Request) {
	p := a.page(typeTags)
	p.Title = "标签"
	p.Canonical = vars.TagsURL()

	p.render(w, "tags", nil)
}

// 主题文件
// /themes/...
func (a *app) getThemes(w http.ResponseWriter, r *http.Request) {
	if isIgnoreThemeFile(r.URL.Path) { // 不展示模板文件
		a.renderError(w, http.StatusNotFound)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, vars.ThemesURL(""))

	if len(path) < len(r.URL.Path) {
		filename := filepath.Join(a.path.ThemesDir, path)

		if !utils.FileExists(filename) {
			a.renderError(w, http.StatusNotFound)
			return
		}

		stat, err := os.Stat(filename)
		if err != nil {
			logs.Error(err)
			a.renderError(w, http.StatusInternalServerError)
			return
		}

		if stat.IsDir() {
			a.renderError(w, http.StatusForbidden)
			return
		}

		http.ServeFile(w, r, filename)
		return
	}

	a.renderError(w, http.StatusNotFound)
	return
}

// /search.html?q=key&page=2
func (a *app) getSearch(w http.ResponseWriter, r *http.Request) {
	p := a.page(typeSearch)

	q := r.FormValue("q")
	if len(q) == 0 {
		http.Redirect(w, r, vars.PostsURL(1), http.StatusPermanentRedirect)
		return
	}

	page, ok := a.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("参数 page: %d 小于 1", page)
		a.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	// 查找标题和内容
	posts := make([]*data.Post, 0, len(a.buf.Data.Posts))
	key := strings.ToLower(q)
	for _, v := range a.buf.Data.Posts {
		if strings.Contains(v.Title, key) || strings.Contains(v.Content, key) {
			posts = append(posts, v)
		}
	}

	p.Title = "搜索:" + q
	p.Q = q
	p.Keywords = q + ",搜索,search"
	p.Description = "搜索关键字" + q + "的结果"
	start, end, ok := a.getPostsRange(len(posts), page, w)
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

// 读取根下的文件
// /...
func (a *app) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		a.getPosts(w, r)
		return
	}

	if !utils.FileExists(filepath.Join(a.path.RawsDir, r.URL.Path)) {
		a.renderError(w, http.StatusNotFound)
		return
	}

	prefix := "/"
	root := http.Dir(a.path.RawsDir)
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// 确认当前文章列表页选择范围。
func (a *app) getPostsRange(postsSize, page int, w http.ResponseWriter) (start, end int, ok bool) {
	size := a.buf.Data.Config.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%d]，实际文章数量为[%d]\n", page, postsSize)
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
		logs.Infof("%s: %s", r.UserAgent(), r.URL) // 输出访问日志

		// 直接根据整个博客的最后更新时间来确认 etag
		if r.Header.Get("If-None-Match") == a.buf.Etag {
			logs.Infof("304: %s", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", a.buf.Etag)
		w.Header().Set("Content-Language", a.buf.Data.Config.Language)
		compress.New(f, logs.ERROR()).ServeHTTP(w, r)
	}
}

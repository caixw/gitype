// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

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

func (c *Client) removeRoutes() {
	for _, route := range c.routes {
		c.mux.Remove(route)
	}

	c.routes = nil
}

func (c *Client) initRoutes() {
	// posts/2016/about.html
	pattern := vars.Post + "/{slug}" + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getPost))

	// index.html
	pattern = vars.Posts + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getPosts))

	// tags/tag1.html
	pattern = vars.Tag + "/{slug}" + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getTag))

	// tags.html
	pattern = vars.Tags + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getTags))

	// search.html
	pattern = vars.Search + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getSearch))

	// themes/...
	pattern = vars.Themes + "/{path}"
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getThemes))

	// /...
	pattern = "/{path}"
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getRaws))
}

// 文章详细页
// /posts/{slug}.html
func (c *Client) getPost(w http.ResponseWriter, r *http.Request) {
	id, found := c.paramString(w, r, "slug")
	if !found {
		return
	}

	var post *data.Post
	var next, prev *data.Link
	for index, p := range c.data.Posts {
		if p.Slug != id {
			continue
		}
		post = p

		if index > 0 {
			p := c.data.Posts[index-1]
			prev = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}

		index++
		if index < len(c.data.Posts) {
			p := c.data.Posts[index]
			next = &data.Link{
				Text: p.Title,
				URL:  p.Permalink,
			}
		}
	} // end for a.data.Posts

	if post == nil {
		logs.Debugf("并未找到与之相对应的文章:%v", id)
		c.getRaws(w, r) // 文章不存在，则查找raws目录下是否存在同名文件
		return
	}

	p := c.newPage()
	p.Post = post
	p.NextPage = next
	p.PrevPage = prev
	p.Canonical = post.Permalink
	p.Description = post.Summary
	p.Title = post.Title
	p.Keywords = post.Keywords

	p.render(w, post.Template, nil)
}

// 首页及文章列表页
// /
// /posts.html?page=2
func (c *Client) getPosts(w http.ResponseWriter, r *http.Request) {
	page, ok := c.queryInt(w, r, "page", 1)
	if !ok {
		return
	}

	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1\n", page)
		c.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到404页面
		return
	}

	p := c.newPage()
	if page > 1 { // 非首页，标题显示页码数
		p.Title = fmt.Sprintf("第%v页", page)
	}
	p.Canonical = c.postsURL(page)

	start, end, ok := c.getPostsRange(len(c.data.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = c.data.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  c.postsURL(page - 1), // 页码从1开始计数
		}
	}
	if end < len(c.data.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  c.postsURL(page + 1),
		}
	}

	p.render(w, "posts", nil)
}

// 标签详细页
// /tags/tag1.html?page=2
func (c *Client) getTag(w http.ResponseWriter, r *http.Request) {
	slug, ok := c.paramString(w, r, "slug")
	if !ok {
		return
	}

	var tag *data.Tag
	for _, t := range c.data.Tags {
		if t.Slug == slug {
			tag = t
			break
		}
	}

	if tag == nil {
		logs.Debugf("查找的标签[%v]不存在", slug)
		c.getRaws(w, r) // 标签不存在，则查找该文件是否存在于 raws 目录下。
		return
	}

	page, ok := c.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1", page)
		c.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := c.newPage()
	p.Tag = tag
	p.Title = tag.Title
	p.Canonical = c.tagURL(slug, page)
	p.Description = "标签" + tag.Title + "的介绍"

	start, end, ok := c.getPostsRange(len(tag.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = tag.Posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  c.tagURL(slug, page-1), // 页码从1开始计数
		}
	}
	if end < len(tag.Posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  c.tagURL(slug, page+1), // 页码从1开始计数
		}
	}

	p.render(w, "tag", nil)
}

// 标签列表页
// /tags.html
func (c *Client) getTags(w http.ResponseWriter, r *http.Request) {
	p := c.newPage()
	p.Title = "标签"
	p.Canonical = vars.Tags + vars.Suffix

	p.render(w, "tags", nil)
}

// 主题文件
// /themes/...
func (c *Client) getThemes(w http.ResponseWriter, r *http.Request) {
	root := http.Dir(c.path.ThemesDir)
	http.StripPrefix(vars.Themes, http.FileServer(root)).ServeHTTP(w, r)
}

// /search.html?q=key&page=2
func (c *Client) getSearch(w http.ResponseWriter, r *http.Request) {
	p := c.newPage()

	key := r.FormValue("q")
	if len(key) == 0 {
		p.render(w, "search", nil)
		return
	}

	page, ok := c.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("参数 page: %v 小于 1", page)
		c.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到404页面
		return
	}

	// 查找标题和内容
	posts := make([]*data.Post, 0, c.data.Config.PageSize)
	for _, v := range c.data.Posts {
		if strings.Index(v.Title, key) >= 0 || strings.Index(v.Content, key) >= 0 {
			posts = append(posts, v)
		}
	}

	p.Title = "搜索:" + key
	p.Keywords = key
	p.Description = "搜索关键字" + key + "的结果集"
	start, end, ok := c.getPostsRange(len(posts), page, w)
	if !ok {
		return
	}
	p.Posts = posts[start:end]
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  c.searchURL(key, page-1), // 页码从1开始计数
		}
	}
	if end < len(posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  c.searchURL(key, page+1),
		}
	}

	p.render(w, "search", nil)

}

// 读取根下的文件
// /...
func (c *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		c.getPosts(w, r)
		return
	}

	root := http.Dir(c.path.RawsDir)
	if !utils.FileExists(filepath.Join(c.path.RawsDir, r.URL.Path)) {
		c.renderError(w, http.StatusNotFound)
		return
	}
	prefix := "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)

}

// 确认当前文章列表页选择范围。
func (c *Client) getPostsRange(postsSize, page int, w http.ResponseWriter) (start, end int, ok bool) {
	size := c.data.Config.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%v]，实际文章数量为[%v]\n", page, postsSize)
		c.renderError(w, http.StatusNotFound) // 页码超出范围，不存在
		return 0, 0, false
	}

	end = start + size
	if postsSize < end {
		end = postsSize
	}

	return start, end, true
}

// 每次访问前需要做的预处理工作。
func (c *Client) pre(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 输出访问日志
		logs.Infof("%v：%v", r.UserAgent(), r.URL)

		// 直接根据整个博客的最后更新时间来确认etag
		if r.Header.Get("If-None-Match") == c.etag {
			logs.Infof("304:%v", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", c.etag)
		compress.New(f).ServeHTTP(w, r)
	}
}

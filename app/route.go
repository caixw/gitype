// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"

	"github.com/caixw/typing/data"
	"github.com/issue9/logs"
	"github.com/issue9/web"
)

func (a *App) initRoute() error {
	m, err := web.NewModule("front")
	if err != nil {
		return err
	}

	urls := a.data.URLS
	p := m.Prefix(urls.Root)

	p.GetFunc(urls.Post+"/{slug}"+urls.Suffix, accessLog(a.getPost)).
		GetFunc(urls.Posts+urls.Suffix, accessLog(a.getPosts)).
		GetFunc(urls.Tag+"/{slug}"+urls.Suffix, accessLog(a.getTag)).
		GetFunc(urls.Tags+urls.Suffix+"{:.*}", accessLog(a.getTags)).
		GetFunc(urls.Themes+"/", accessLog(a.getThemes)).
		GetFunc("/", accessLog(a.getRaws))

	// feeds
	conf := a.data.Config
	if conf.RSS != nil {
		p.GetFunc(conf.RSS.URL, accessLog(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.rssBuffer.Bytes())
		}))
	}

	if conf.Atom != nil {
		p.GetFunc(conf.Atom.URL, accessLog(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.atomBuffer.Bytes())
		}))
	}

	if conf.Sitemap != nil {
		p.GetFunc(conf.Sitemap.URL, accessLog(func(w http.ResponseWriter, r *http.Request) {
			w.Write(a.sitemapBuffer.Bytes())
		}))
	}
	return nil
}

// /
func (a *App) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == a.data.URLS.Root || r.URL.Path == a.data.URLS.Root+"/" {
		a.getPosts(w, r)
		return
	}

	root := http.Dir(a.path.DataRaws)
	prefix := a.data.URLS.Root + "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// 从posts中摘取指定页码的文章存入到p中。
// posts用于筛选的所有文章列表；page当前显示的页码，从1开始。
func (a *App) getPagePost(p *page, posts []*data.Post, page int, w http.ResponseWriter) (ok bool) {
	size := a.data.Config.PageSize
	start := size * (page - 1) // 系统从零开始计数
	if start > len(posts) {
		logs.Debugf("请示页码为[%v]，实际文章数量为[%v]", page, len(posts))
		w.WriteHeader(http.StatusNotFound) // 页码超出范围，不存在
		return false
	}

	end := start + size
	if len(posts) < end {
		end = len(posts)
	}

	p.Posts = posts[start:end]
	p.Canonical = a.postsURL(uint(page))
	if page > 1 {
		p.PrevPage = &data.Link{
			Text: "前一页",
			URL:  a.postsURL(uint(page - 1)), // 页码从1开始计数
		}
	}
	if end < len(posts) {
		p.PrevPage = &data.Link{
			Text: "下一页",
			URL:  a.postsURL(uint(page + 1)), // 页码从1开始计数
		}
	}

	return true
}

// 首页及文章列表页
// /
// /posts.html?page=2
func (a *App) getPosts(w http.ResponseWriter, r *http.Request) {
	page, ok := queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		w.WriteHeader(http.StatusNotFound) // 页码为负数的页码不存在
		return
	}

	p := a.newPage()
	if !a.getPagePost(p, a.data.Posts, page, w) {
		return
	}
	p.render(w, r, "posts", nil)
}

// 主题文件
// /themes/...
func (a *App) getThemes(w http.ResponseWriter, r *http.Request) {
	root := http.Dir(a.path.DataThemes)
	prefix := a.data.URLS.Root + a.data.URLS.Themes
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// 标签详细页
// /tags/tag1.html?page=2
func (a *App) getTag(w http.ResponseWriter, r *http.Request) {
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
		logs.Debugf("查找的标签[%v]不存在", slug)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page, ok := queryInt(w, r, "page", 1)
	if !ok {
		return
	}

	if page < 1 {
		logs.Debugf("请求的页码[%v]小于1", page)
		w.WriteHeader(http.StatusNotFound) // 页码为负数的页码不存在
		return
	}

	p := a.newPage()
	p.Tag = tag
	if !a.getPagePost(p, tag.Posts, page, w) {
		return
	}
	p.render(w, r, "tag", nil)
}

// 标签列表页
// /tags.html
func (a *App) getTags(w http.ResponseWriter, r *http.Request) {
	p := a.newPage()
	p.Title = "标签"
	p.Canonical = a.tagsURL()
	p.render(w, r, "tags", nil)
}

// 文章详细页
// /posts/{slug}.html
func (a *App) getPost(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	p := a.newPage()
	p.Post = post
	p.NextPage = next
	p.PrevPage = prev
	p.Canonical = post.Permalink

	p.render(w, r, post.Template, nil)
}

// 输出访问日志
func accessLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs.Infof("%v：%v", r.UserAgent(), r.URL)
		h(w, r)
	}
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/middleware/compress"
	"github.com/issue9/mux"
	"github.com/issue9/mux/params"
	"github.com/issue9/utils"
)

// 模板的扩展名，在主题目录下，以下扩展名的文件，不会被展示
var ignoreThemeFileExts = []string{
	vars.TemplateExtension,
	".yaml",
	".yml",
}

func isIgnoreThemeFile(file string) bool {
	ext := filepath.Ext(file)

	for _, v := range ignoreThemeFileExts {
		if ext == v {
			return true
		}
	}

	return false
}

func (client *Client) initRoutes() (err error) {
	handle := func(pattern string, h http.HandlerFunc) {
		if err != nil {
			return
		}

		client.patterns = append(client.patterns, pattern)
		err = client.mux.HandleFunc(pattern, client.prepare(h), http.MethodGet)
	}

	handle(vars.PostURL("{slug}"), client.getPost)     // posts/2016/about.html   posts/{slug}.html
	handle(vars.IndexURL(0), client.getPosts)          // index.html
	handle(vars.LinksURL(), client.getLinks)           // links.html
	handle(vars.TagURL("{slug}", 1), client.getTag)    // tags/tag1.html     tags/{slug}.html
	handle(vars.TagsURL(), client.getTags)             // tags.html
	handle(vars.ArchivesURL(), client.getArchives)     // archives.html
	handle(vars.SearchURL("", 1), client.getSearch)    // search.html
	handle(vars.ThemesURL("{path}"), client.getThemes) // themes/...          themes/{path}
	handle("/{path}", client.getRaws)                  // /...                /{path}

	return err
}

// 文章详细页
// /posts/{slug}.html
func (client *Client) getPost(w http.ResponseWriter, r *http.Request) {
	id, found := client.paramString(w, r, "slug")
	if !found {
		return
	}

	var index int
	for i, p := range client.data.Posts {
		if p.Slug == id {
			index = i
			break
		}
	}

	if index < 0 {
		logs.Debugf("并未找到与之相对应的文章:%s", id)
		client.getRaws(w, r) // 文章不存在，则查找 raws 目录下是否存在同名文件
		return
	}

	post := client.data.Posts[index]
	p := client.page(typePost)

	client.data.Outdated(post)
	p.Post = post
	p.Keywords = post.Keywords
	p.Description = post.Summary
	p.Title = post.Title
	p.Canonical = client.data.URL(post.Permalink)
	p.License = post.License // 文章可具体指定协议
	p.Author = post.Author   // 文章可具体指定作者

	if index > 0 {
		prev := client.data.Posts[index-1]
		p.prevPage(prev.Permalink, prev.Title)
	}
	if index+1 < len(client.data.Posts) {
		next := client.data.Posts[index+1]
		p.nextPage(next.Permalink, next.Title)
	}

	p.render(w, post.Template, nil)
}

// 首页及文章列表页
// /
// /posts.html?page=2
func (client *Client) getPosts(w http.ResponseWriter, r *http.Request) {
	page, ok := client.queryInt(w, r, "page", 1)
	if !ok {
		return
	}

	if page < 1 {
		logs.Debugf("请求的页码[%d]小于1\n", page)
		client.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := client.page(typeIndex)
	if page > 1 { // 非首页，标题显示页码数
		p.Type = typePosts
		p.Title = fmt.Sprintf("第 %d 页", page)
	}
	p.Canonical = client.data.URL(vars.PostsURL(page))

	start, end, ok := client.getPostsRange(len(client.data.Posts), page, w)
	if !ok {
		return
	}
	p.Posts = client.data.Posts[start:end]
	if page > 1 {
		p.prevPage(vars.PostsURL(page-1), "")
	}
	if end < len(client.data.Posts) {
		p.nextPage(vars.PostsURL(page+1), "")
	}

	p.render(w, "posts", nil)
}

// 标签详细页
// /tags/tag1.html?page=2
func (client *Client) getTag(w http.ResponseWriter, r *http.Request) {
	slug, ok := client.paramString(w, r, "slug")
	if !ok {
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
		client.getRaws(w, r) // 标签不存在，则查找该文件是否存在于 raws 目录下。
		return
	}

	page, ok := client.queryInt(w, r, "page", 1)
	if !ok {
		return
	}
	if page < 1 {
		logs.Debugf("请求的页码[%d]小于1", page)
		client.renderError(w, http.StatusNotFound) // 页码为负数的表示不存在，跳转到 404 页面
		return
	}

	p := client.page(typeTag)
	p.Tag = tag
	p.Title = tag.Title
	p.Keywords = tag.Keywords
	p.Description = tag.Description
	p.Canonical = client.data.URL(vars.TagURL(slug, page))

	start, end, ok := client.getPostsRange(len(tag.Posts), page, w)
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
func (client *Client) getLinks(w http.ResponseWriter, r *http.Request) {
	p := client.page(typeLinks)
	p.Title = "友情链接"
	p.Canonical = client.data.URL(vars.LinksURL())

	p.render(w, "links", nil)
}

// 标签列表页
// /tags.html
func (client *Client) getTags(w http.ResponseWriter, r *http.Request) {
	p := client.page(typeTags)
	p.Title = "标签"
	p.Canonical = client.data.URL(vars.TagsURL())
	p.Description = "标签列表"

	p.render(w, "tags", nil)
}

// 归档页
// /archives.html
func (client *Client) getArchives(w http.ResponseWriter, r *http.Request) {
	p := client.page(typeArchives)
	p.Title = "归档"
	p.Keywords = "归档,存档,archive,archives"
	p.Description = "网站的归档列表，按时间进行排序"
	p.Canonical = client.data.URL(vars.ArchivesURL())
	p.Archives = client.data.Archives

	p.render(w, "archives", nil)
}

// 主题文件
// /themes/...
func (client *Client) getThemes(w http.ResponseWriter, r *http.Request) {
	if isIgnoreThemeFile(r.URL.Path) { // 不展示模板文件，查看 raws 中是否有同名文件
		client.getRaws(w, r)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, vars.ThemesURL(""))
	if len(path) >= len(r.URL.Path) { // path 不包含 vars.ThemesURL("") 前缀
		client.getRaws(w, r)
		return
	}

	filename := filepath.Join(client.path.ThemesDir, path)

	if !utils.FileExists(filename) {
		client.getRaws(w, r)
		return
	}

	stat, err := os.Stat(filename)
	if err != nil {
		logs.Error(err)
		client.renderError(w, http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		client.getRaws(w, r)
		return
	}

	http.ServeFile(w, r, filename)
}

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

	// 查找标题和内容
	posts := make([]*data.Post, 0, len(client.data.Posts))
	key := strings.ToLower(q)
	for _, v := range client.data.Posts {
		if strings.Contains(v.Title, key) || strings.Contains(v.Content, key) {
			posts = append(posts, v)
		}
	}

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

// /...
func (client *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		client.getPosts(w, r)
		return
	}

	if !utils.FileExists(filepath.Join(client.path.RawsDir, r.URL.Path)) {
		client.renderError(w, http.StatusNotFound)
		return
	}

	prefix := "/"
	root := http.Dir(client.path.RawsDir)
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// 确认当前文章列表页选择范围。
func (client *Client) getPostsRange(postsSize, page int, w http.ResponseWriter) (start, end int, ok bool) {
	size := client.data.Config.PageSize
	start = size * (page - 1) // 系统从零开始计数
	if start > postsSize {
		logs.Debugf("请求页码为[%d]，实际文章数量为[%d]\n", page, postsSize)
		client.renderError(w, http.StatusNotFound) // 页码超出范围，不存在
		return 0, 0, false
	}

	end = start + size
	if postsSize < end {
		end = postsSize
	}

	return start, end, true
}

// 每次访问前需要做的预处理工作。
func (client *Client) prepare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs.Infof("%s: %s", r.UserAgent(), r.URL) // 输出访问日志

		// 直接根据整个博客的最后更新时间来确认 etag
		if r.Header.Get("If-None-Match") == client.etag {
			logs.Infof("304: %s", r.URL)
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Etag", client.etag)
		w.Header().Set("Content-Language", client.data.Config.Language)
		compress.New(f, logs.ERROR()).ServeHTTP(w, r)
	}
}

// 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回 false
func (client *Client) paramString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	ps := mux.Params(r)
	val, err := ps.String(key)

	if err == params.ErrParamNotExists {
		client.renderError(w, http.StatusNotFound)
		return "", false
	} else if err != nil {
		logs.Error(err)
		client.renderError(w, http.StatusNotFound)
		return "", false
	} else if len(val) == 0 {
		client.renderError(w, http.StatusNotFound)
		return "", false
	}

	return val, true
}

// 获取查询参数 key 的值，并将其转换成 Int 类型，若该值不存在返回 def 作为其默认值，
// 若是类型不正确，则返回一个 false，并向客户端输出一个 400 错误。
func (client *Client) queryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		logs.Error(err)
		client.renderError(w, http.StatusBadRequest)
		return 0, false
	}
	return ret, true
}

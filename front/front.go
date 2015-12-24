// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"database/sql"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/util"
	"github.com/issue9/conv"
	"github.com/issue9/handlers"
	"github.com/issue9/is"
	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/orm/fetch"
	"github.com/issue9/web"
)

var (
	cfg  *app.Config
	opt  *app.Options
	stat *app.Stat
	db   *orm.DB
)

// 从主题根目录加载所有的主题内容，并初始所有的主题下静态文件的路由。
// defaultTheme 为默认的主题。
func Init(c *app.Config, database *orm.DB, options *app.Options, s *app.Stat) error {
	cfg = c
	opt = options
	db = database
	stat = s

	if err := loadThemes(); err != nil {
		return err
	}

	if err := Switch(opt.Theme); err != nil {
		return err
	}

	return initRoute()
}

func initRoute() error {
	m, err := web.NewModule("front")
	if err != nil {
		return err
	}

	m.Get(opt.HomeURL(), etagHandler(handlers.CompressFunc(pageHome))).
		Get(opt.TagsURL(), etagHandler(handlers.CompressFunc(pageTags))).
		Get(opt.TagURL("{id}", 1), etagHandler(handlers.CompressFunc(pageTag))).
		Get(opt.PostsURL(1), etagHandler(handlers.CompressFunc(pagePosts))).
		Get(opt.PostURL("{id}"), etagHandler(handlers.CompressFunc(pagePost))). // 获取文章详细内容
		Post(opt.PostURL("{id}"), etagHandler(handlers.CompressFunc(pagePost))) // 提交评论

	// 静态文件路由，TODO 去掉config中对于必须以/结尾的判断
	// TODO 静态文件压缩
	m.Get(cfg.UploadURLPrefix+"/", http.StripPrefix(cfg.UploadURLPrefix, http.FileServer(http.Dir(cfg.UploadDir)))).
		Get(cfg.ThemeURLPrefix+"/", http.StripPrefix(cfg.ThemeURLPrefix, http.FileServer(http.Dir(cfg.ThemeDir))))

	m.Prefix(cfg.FrontAPIPrefix).
		PostFunc("/posts/{id:\\d+}/comments", frontPostPostComment).
		GetFunc("/posts/{id:\\d+}/comments", frontGetPostComments)

	return nil
}

// etag包装
func etagHandler(h http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		etag := strconv.FormatInt(opt.LastUpdated, 10)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		w.Header().Set("Etag", etag)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}

func getTagPosts(page int, tagID int64) ([]*Post, error) {
	posts := make([]*Post, 0, opt.PageSize)
	sql := `SELECT p.{id} AS {ID}, p.{name} AS {Name}, p.{title} AS {Title}, p.{summary} AS {Summary},
		p.{content} as {Content}, p.{created} AS {Created}, p.{allowComment} AS {AllowComment}
		FROM #relationships AS r
		LEFT JOIN #posts AS p ON p.{id}=r.{postID}
		WHERE p.{state}=? AND r.{tagID}=?
		ORDER BY {order} ASC, {created} DESC
		LIMIT ? OFFSET ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, tagID, opt.PageSize, opt.PageSize*page)
	if err != nil {
		return nil, err
	}
	_, err = fetch.Obj(&posts, rows)
	rows.Close()

	return posts, err
}

func getPosts(page int) ([]*Post, error) {
	posts := make([]*Post, 0, opt.PageSize)
	sql := `SELECT {id} AS {ID}, {name} AS {Name}, {title} AS {Title}, {summary} AS {Summary},
	{content} AS {Content}, {created} AS {Created}, {modified} AS {Mofified}, {allowComment} AS {AllowComment}
	FROM #posts
	WHERE {state}=?
	ORDER BY {order} ASC, {created} DESC
	LIMIT ? OFFSET ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, opt.PageSize, opt.PageSize*page)
	if err != nil {
		return nil, err
	}
	_, err = fetch.Obj(&posts, rows)
	rows.Close()

	return posts, err
}

// 输出一个特写状态码下的错误页面。若该页面模板不存在，则只输出状态码，而没有内容。
// 只对状态码大于等于400的起作用。
func pageHttpStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	if code < 400 {
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)

	path := themeDir(currentTheme) + strconv.Itoa(code) + ".html"
	if !util.FileExists(path) { // 文件不存在，则只输出状态码，省略内容。
		return
	}

	// TODO: serveFile会自动写入一个状态码，导致多次输出状态码的提示
	http.ServeFile(w, r, path)
}

// 首页
func pageHome(w http.ResponseWriter, r *http.Request) {
	// 首页的匹配模式为：/，可以匹配任意路径。所以此处作个判断，除了完全匹配的，其余都返回404
	if r.URL.Path != opt.HomeURL() {
		pageHttpStatusCode(w, r, http.StatusNotFound)
		return
	}

	pagePosts(w, r)
}

// 首页或是列表页
func pagePosts(w http.ResponseWriter, r *http.Request) {
	info, err := getInfo()
	if err != nil {
		logs.Error("pagePosts:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	page := conv.MustInt(r.FormValue("page"), 1)
	if page == 1 {
		info.Canonical = opt.URL(opt.HomeURL())
	} else if page > 1 { // 为1的时候，不需要prev
		info.Canonical = opt.URL(opt.PostsURL(page))
		info.PrevPage = &Anchor{Title: "上一页", Link: opt.PostsURL(page - 1)}
	}

	if page*opt.SidebarSize < info.PostSize {
		info.NextPage = &Anchor{Title: "下一页", Link: opt.PostsURL(page + 1)}
	}

	posts, err := getPosts(page - 1)
	if err != nil {
		logs.Error("pagePosts:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"info":  info,
		"posts": posts,
	}
	render(w, r, "posts", data, map[string]string{"Content-Type": "text/html"})
}

// /tags
func pageTags(w http.ResponseWriter, r *http.Request) {
	info, err := getInfo()
	if err != nil {
		logs.Error("pageTags:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}
	info.Canonical = opt.URL(opt.TagsURL())
	info.Title = "标签"

	sql := `SELECT {id} AS {ID}, {name} AS {Name}, {title} AS {Title} FROM #tags`
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("pageTags:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tags := make([]*Tag, 0, 100)
	if _, err = fetch.Obj(&tags, rows); err != nil {
		logs.Error("pageTags:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"info": info, "tags": tags}
	render(w, r, "tags", data, map[string]string{"Content-Type": "text/html"})
}

// /tags/1.html
func pageTag(w http.ResponseWriter, r *http.Request) {
	tagName, ok := util.ParamString(w, r, "id")
	if !ok {
		return
	}
	tagName = strings.TrimSuffix(tagName, opt.Suffix)

	sql := `SELECT {id} AS {ID}, {name} AS {Name}, {title} AS {Title}, {description} AS {Description}
	FROM #tags
	WHERE {name}=?`
	rows, err := db.Query(true, sql, tagName)
	if err != nil {
		logs.Error("pageTag:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tag := &Tag{}
	if _, err = fetch.Obj(tag, rows); err != nil {
		logs.Error("pageTag:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	info, err := getInfo()
	if err != nil {
		logs.Error("pageTag:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	info.Canonical = opt.URL(tag.Permalink())
	info.Title = tag.Title

	page := conv.MustInt(r.FormValue("page"), 1)
	if page < 1 { // 不能小于1
		page = 1
	} else if page > 1 { // 为1的时候，不需要prev
		info.PrevPage = &Anchor{Title: "上一页", Link: opt.TagURL(tagName, page-1)}
	}
	if page*opt.SidebarSize < tag.Count() {
		info.NextPage = &Anchor{Title: "下一页", Link: opt.TagURL(tagName, page+1)}
	}
	posts, err := getTagPosts(page-1, tag.ID)
	if err != nil {
		logs.Error("pageTag:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"info":  info,
		"tag":   tag,
		"posts": posts,
	}
	render(w, r, "tag", data, map[string]string{"Content-Type": "text/html"})
}

// /posts/1.html
// /posts/about.html
func pagePost(w http.ResponseWriter, r *http.Request) {
	idStr, ok := util.ParamString(w, r, "id")
	if !ok {
		return
	}
	idStr = strings.TrimSuffix(idStr, opt.Suffix)

	var rows *sql.Rows
	var err error
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		sql := `SELECT * FROM #posts WHERE {name}=?`
		rows, err = db.Query(true, sql, idStr)
	} else {
		sql := `SELECT * FROM #posts WHERE {id}=?`
		rows, err = db.Query(true, sql, postID)
	}
	if err != nil {
		logs.Error("pagePost:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	mp := &models.Post{}
	if _, err = fetch.Obj(mp, rows); err != nil {
		logs.Error("pagePost:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}

	if len(mp.Title) == 0 || mp.State != models.PostStatePublished {
		pageHttpStatusCode(w, r, http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		if err := insertComment(mp.ID, r); err != nil {
			logs.Error("pagePost:", err)
		} else {
			stat.WaitingCommentsSize++
			stat.CommentsSize++
		}
	}

	post := &Post{
		ID:           mp.ID,
		Name:         mp.Name,
		Title:        mp.Title,
		Summary:      mp.Summary,
		Content:      mp.Content,
		Author:       opt.ScreenName,
		Created:      mp.Created,
		Modified:     mp.Modified,
		AllowComment: mp.AllowComment,
	}

	info, err := getInfo()
	if err != nil {
		logs.Error("pagePost:", err)
		pageHttpStatusCode(w, r, http.StatusInternalServerError)
		return
	}
	info.Canonical = opt.URL(post.Permalink())
	info.Title = post.Title
	info.Description = post.Summary
	info.Keywords = post.Keywords()

	data := map[string]interface{}{
		"info": info,
		"post": post,
	}
	render(w, r, "post", data, map[string]string{"Content-Type": "text/html"})
}

// 将当前提交的评论插入数据库
func insertComment(postID int64, r *http.Request) error {
	c := &models.Comment{
		//Parent  int64  `orm:"name(parent)"`          // 子评论的话，此为其上一级评论的id
		Created:     time.Now().Unix(),
		PostID:      postID,
		State:       models.CommentStateWaiting,
		IP:          r.RemoteAddr,
		Agent:       r.UserAgent(),
		IsAdmin:     false,
		Content:     r.FormValue("content"),
		AuthorName:  r.FormValue("name"),
		AuthorEmail: r.FormValue("email"),
		AuthorURL:   r.FormValue("url"),
	}
	_, err := db.Insert(c)
	return err
}

// @api get /api/posts/{id}/comments
// @apiQuery page  int 页码
// @apiGroup front
//
// @apiSuccess 200 OK
// @apiParam comments array 当前页的评论
func frontGetPostComments(w http.ResponseWriter, r *http.Request) {
	id, ok := util.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("frontGetPostComments:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != models.PostStatePublished {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	sql := db.Where("{postID}=?", id).
		And("{state}=?", models.CommentStateApproved).
		Table("#comments")

	var page int
	if page, ok = util.QueryInt(w, r, "page", 0); !ok {
		return
	}
	sql.Limit(opt.PageSize, page*opt.PageSize)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("frontGetComments:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, map[string]interface{}{"comments": maps}, nil)
}

// @api post /api/posts/{id}/comments 提交新评论
// @apiGroup front
//
// @apiRequest json
// @apiParam parent      int    评论的父级内容
// @apiParam postID      int    评论的文章
// @apiParam content     string 评论的内容
// @apiParam authorName  string 评论的作者
// @apiParam authorURL   string 评论作者的网站地址，可为空
// @apiParam authorEmail string 评论作者的邮箱
//
// @apiSuccess 201 created
func frontPostPostComment(w http.ResponseWriter, r *http.Request) {
	c := &struct {
		Parent      int64  `json:"parent"`
		PostID      int64  `json:"postID"`
		Content     string `json:"content"`
		AuthorName  string `json:"authorName"`
		AuthorURL   string `json:"authorURL"`
		AuthorEmail string `json:"authorEmail"`
	}{}

	if !util.ReadJSON(w, r, c) {
		return
	}

	// 判断文章状态
	if c.PostID <= 0 {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	p := &models.Post{ID: c.PostID}
	if err := db.Select(p); err != nil {
		logs.Error("forntPostPostComment:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if (len(p.Title) == 0 && len(p.Content) == 0) || p.State != models.PostStatePublished {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}
	if !p.AllowComment {
		util.RenderJSON(w, http.StatusMethodNotAllowed, nil, nil)
		return
	}

	// 判断提交数据的状态
	errs := &util.ErrorResult{}
	if c.Parent < 0 {
		errs.Detail["parent"] = "无效的parent"
	}
	if len(c.Content) == 0 {
		errs.Detail["content"] = "content不能为空"
	}
	if len(c.AuthorURL) > 0 && !is.URL(c.AuthorURL) {
		errs.Detail["authorURL"] = "无效的authorURL"
	}
	if !is.Email(c.AuthorEmail) {
		errs.Detail["authorEmail"] = "无效的authorEmail"
	}
	if len(c.AuthorName) == 0 {
		errs.Detail["authorName"] = "authorName不能为空"
	}

	c.AuthorName = html.EscapeString(c.AuthorName)

	// url只提取其host部分，其余的都去掉
	u, err := url.Parse(c.AuthorURL)
	if err != nil {
		logs.Error("frontPostComment:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	c.AuthorURL = u.Scheme + ":" + u.Host

	c.Content = html.EscapeString(c.Content)
	c.Content = strings.Replace(c.Content, "\n", "<br />", -1)

	comm := &models.Comment{
		PostID:      c.PostID,
		Parent:      c.Parent,
		AuthorURL:   c.AuthorURL,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Content:     c.Content,
		Created:     time.Now().Unix(),
		State:       models.CommentStateWaiting,
		IP:          r.RemoteAddr,
		Agent:       r.UserAgent(),
		IsAdmin:     false,
	}
	if _, err := db.Insert(comm); err != nil {
		logs.Error("frontPostComment:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	util.RenderJSON(w, http.StatusCreated, nil, nil)
}

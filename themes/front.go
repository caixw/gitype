// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"database/sql"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/conv"
	"github.com/issue9/is"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
	"github.com/issue9/web"
)

func initRoute() error {
	m, err := web.NewModule("front")
	if err != nil {
		return err
	}

	m.GetFunc("/", pagePosts).
		GetFunc("/tags"+opt.Suffix, pageTags).
		GetFunc("/tags/{id}"+opt.Suffix, pageTag).
		GetFunc("/posts"+opt.Suffix, pagePosts).
		GetFunc("/posts/{id}"+opt.Suffix, pagePost).  // 获取文章详细内容
		PostFunc("/posts/{id}"+opt.Suffix, pagePost). // 提交评论
		Get(cfg.UploadURLPrefix, http.StripPrefix(cfg.ThemeURLPrefix, http.FileServer(http.Dir(cfg.UploadDir)))).
		Get(cfg.ThemeURLPrefix, http.StripPrefix(cfg.ThemeURLPrefix, http.FileServer(http.Dir(cfg.ThemeDir))))

	m.Prefix(cfg.FrontAPIPrefix).
		PostFunc("/posts/{id:\\d+}/comments", frontPostPostComment).
		GetFunc("/posts/{id:\\d+}/comments", frontGetPostComments)

	return nil
}

func getTagPosts(page int, tagID int64) ([]*Post, error) {
	posts := make([]*Post, 0, opt.PageSize)
	sql := `SELECT p.{id} AS ID, p.{name} AS Name, p.{title} AS Title, p.{summary} AS Summary,
		p.{content} as Content, p.{created} AS Created, p.{allowComment} AS AllowComment
		FROM #relationships AS r
		LEFT JOIN #posts AS p ON p.{id}=r.{postID}
		WHERE p.{state}=? AND r.{tagID}=?
		ORDER BY {order} DESC
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
	sql := `SELECT {id} AS ID, {name} AS Name, {title} AS Title, {summary} AS Summary,
	{content} AS Content, {created} AS Created, {allowComment} AS AllowComment
	FROM #posts
	WHERE {state}=?
	ORDER BY {order} DESC
	LIMIT ? OFFSET ?`
	rows, err := db.Query(true, sql, models.PostStatePublished, opt.PageSize, opt.PageSize*page)
	if err != nil {
		return nil, err
	}
	_, err = fetch.Obj(&posts, rows)
	rows.Close()

	return posts, err
}

// 首页或是列表页
func pagePosts(w http.ResponseWriter, r *http.Request) {
	info, err := getInfo()
	if err != nil {
		logs.Error("pagePosts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := conv.MustInt(r.FormValue("page"), 1)
	if page < 1 { // 不能小于1
		page = 1
	} else if page > 1 { // 为1的时候，不需要prev
		info.PrevPage = &Anchor{Title: "上一页", Link: core.PostsURL(page - 1)}
	}
	if page*opt.SidebarSize < info.PostSize {
		info.NextPage = &Anchor{Title: "下一页", Link: core.PostsURL(page + 1)}
	}

	posts, err := getPosts(page - 1)
	if err != nil {
		logs.Error("pagePosts:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"info":  info,
		"posts": posts,
	}
	render(w, "posts", data)
}

// /tags
func pageTags(w http.ResponseWriter, r *http.Request) {
	info, err := getInfo()
	if err != nil {
		logs.Error("pageTags:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info.Canonical = opt.SiteURL + "tags"
	info.Title = "标签"

	sql := `SELECT {id} AS ID, {name} AS Name, {title} AS Title FROM #tags`
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("pageTags:", err)
		w.WriteHeader(500)
		return
	}
	defer rows.Close()

	tags := make([]*Tag, 0, 100)
	if _, err = fetch.Obj(&tags, rows); err != nil {
		logs.Error("pageTags:", err)
		w.WriteHeader(500)
		return
	}

	for _, tag := range tags {
		sql := db.Where("tagID=?", tag.ID).Table("#relationships")
		tag.Count, err = sql.Count(true)
		if err != nil {
			logs.Error("pageTags:", err)
			w.WriteHeader(500)
			return
		}
	}
	render(w, "tags", map[string]interface{}{"info": info, "tags": tags})
}

// /tags/1.html
func pageTag(w http.ResponseWriter, r *http.Request) {
	tagName, ok := core.ParamString(w, r, "id")
	if !ok {
		return
	}
	tagName = strings.TrimSuffix(tagName, opt.Suffix)

	sql := `SELECT t.{id} AS ID, t.{name} AS Name, t.{title} AS Title, t.{description} AS Description,
	count(r.{tagID}) AS {Count}
	FROM #tags AS t
	LEFT JOIN #relationships AS r ON t.{id}=r.{tagID}
	WHERE t.{name}=?`
	rows, err := db.Query(true, sql, tagName)
	if err != nil {
		logs.Error("pageTag:", err)
		w.WriteHeader(500)
		return
	}
	defer rows.Close()

	tag := &Tag{}
	if _, err = fetch.Obj(tag, rows); err != nil {
		logs.Error("pageTag:", err)
		w.WriteHeader(500)
		return
	}

	info, err := getInfo()
	if err != nil {
		logs.Error("pageTag:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	info.Canonical = tag.Permalink()
	info.Title = tag.Title

	page := conv.MustInt(r.FormValue("page"), 1)
	if page < 1 { // 不能小于1
		page = 1
	} else if page > 1 { // 为1的时候，不需要prev
		info.PrevPage = &Anchor{Title: "上一页", Link: core.TagURL(tagName, page-1)}
	}
	if page*opt.SidebarSize < tag.Count {
		info.NextPage = &Anchor{Title: "下一页", Link: core.TagURL(tagName, page+1)}
	}
	posts, err := getTagPosts(page-1, tag.ID)
	if err != nil {
		logs.Error("pageTag:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"info":  info,
		"tag":   tag,
		"posts": posts,
	}
	render(w, "tag", data)
}

// /posts/1.html
// /posts/about.html
func pagePost(w http.ResponseWriter, r *http.Request) {
	idStr, ok := core.ParamString(w, r, "id")
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	mp := &models.Post{}
	if _, err = fetch.Obj(mp, rows); err != nil {
		logs.Error("pagePost:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(mp.Title) == 0 || mp.State != models.PostStatePublished {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		if err := insertComment(mp.ID, r); err != nil {
			logs.Error("pagePost:", err)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info.Canonical = post.Permalink()
	info.Title = post.Title

	data := map[string]interface{}{
		"info": info,
		"post": post,
	}
	render(w, "post", data)
}

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
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("frontGetPostComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	sql := db.Where("{postID}=?", id).
		And("{state}=?", models.CommentStateApproved).
		Table("#comments")

	var page int
	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	sql.Limit(opt.PageSize, page*opt.PageSize)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("frontGetComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"comments": maps}, nil)
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

	if !core.ReadJSON(w, r, c) {
		return
	}

	// 判断文章状态
	if c.PostID <= 0 {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	p := &models.Post{ID: c.PostID}
	if err := db.Select(p); err != nil {
		logs.Error("forntPostPostComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if (len(p.Title) == 0 && len(p.Content) == 0) || p.State != models.PostStatePublished {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}
	if !p.AllowComment {
		core.RenderJSON(w, http.StatusMethodNotAllowed, nil, nil)
		return
	}

	// 判断提交数据的状态
	errs := &core.ErrorResult{}
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
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
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
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, nil, nil)
}

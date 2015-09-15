// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"html"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caixw/typing/core"
	"github.com/issue9/is"
	"github.com/issue9/logs"
)

const (
	postStateAll    = iota // 表示所有状态
	postStateNormal        // 正常状态
	postStateDraft         // 草稿
)

// 文章内容
type post struct {
	ID       int64  `orm:"name(id);ai"`
	Name     string `orm:"name(name);len(50);nullable;unique(unq_name)"` // 唯一名称
	Title    string `orm:"name(title);len(200)"`                         // 标题
	Summary  string `orm:"name(summary);len(2000)"`                      // 内容摘要
	Content  string `orm:"name(content);len(-1)"`                        // 实际内容
	State    int    `orm:"name(state)"`                                  // 状态
	Order    int    `orm:"name(order)"`                                  // 排序
	Template string `orm:"name(template)"`                               // 使用的模板
	Password string `orm:"name(password)"`                               // 访问密码

	Created  int64 `orm:"name(created)"`  // 创建时间
	Modified int64 `orm:"name(modified)"` // 最后次修改时间

	AllowPing    bool `orm:"name(allowPing)"`
	AllowComment bool `orm:"name(allowComment)"`
}

func (p *post) Meta() string {
	return `orm:"name(posts)"`
}

// @api post /admin/api/posts 新建文章
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称，可以为空
// @apiParam title string 标题
// @apiParam summary string 文章摘要
// @apiParam content string 文章内容
// @apiParam state int 状态
// @apiParam order int 排序
// @apiParam template string 所使用的模板
// @apiParam password string 访问密码
// @apiParam allowPing bool 允许ping
// @apiParam allowComment bool 允许评论
// @apiParam tags array 关联的标签
// @apiParam cats array 关联的分类
//
// @apiSuccess 201 created
func adminPostPost(w http.ResponseWriter, r *http.Request) {
	p := &struct {
		Name         string  `json:"name"`
		Title        string  `json:"title"`
		Summary      string  `json:"summary"`
		Content      string  `json:"content"`
		State        int     `json:"state"`
		Order        int     `json:"order"`
		Template     string  `json:"template"`
		Password     string  `json:"password"`
		AllowPing    bool    `json:"allowPing"`
		AllowComment bool    `json:"allowComment"`
		Tags         []int64 `json:"tags"`
		Cats         []int64 `json:"cats"`
	}{}

	if !core.ReadJSON(w, r, p) {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("postPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	pp := &post{
		Name:         p.Name,
		Title:        p.Title,
		Summary:      p.Summary,
		Content:      p.Content,
		State:        p.State,
		Order:        p.Order,
		Template:     p.Template,
		Password:     p.Password,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Created:      time.Now().Unix(),
	}

	// 插入文章
	result, err := tx.Insert(pp)
	if err != nil {
		logs.Error("postPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	postID, err := result.LastInsertId()
	if err != nil {
		logs.Error("postPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 插入relationship
	rs := make([]*relationship, 0, len(p.Tags)+len(p.Cats))
	for _, v := range p.Tags {
		rs = append(rs, &relationship{PostID: postID, MetaID: v})
	}
	for _, v := range p.Cats {
		rs = append(rs, &relationship{PostID: postID, MetaID: v})
	}
	if err := tx.InsertMany(rs); err != nil {
		logs.Error("postPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// commit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logs.Error("postPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, nil, nil)
}

// @api put /admin/api/posts/{id} 修改文章
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称，可以为空
// @apiParam title string 标题
// @apiParam summary string 文章摘要
// @apiParam content string 文章内容
// @apiParam state int 状态
// @apiParam order int 排序
// @apiParam template string 所使用的模板
// @apiParam password string 访问密码
// @apiParam allowPing bool 允许ping
// @apiParam allowComment bool 允许评论
// @apiParam tags array 关联的标签
// @apiParam cats array 关联的分类
//
// @apiSuccess 200 no content
func adminPutPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &struct {
		Name         string  `json:"name"`
		Title        string  `json:"title"`
		Summary      string  `json:"summary"`
		Content      string  `json:"content"`
		State        int     `json:"state"`
		Order        int     `json:"order"`
		Template     string  `json:"template"`
		Password     string  `json:"password"`
		AllowPing    bool    `json:"allowPing"`
		AllowComment bool    `json:"allowComment"`
		Tags         []int64 `json:"tags"`
		Cats         []int64 `json:"cats"`
	}{}
	if !core.ReadJSON(w, r, p) {
		return
	}

	pp := &post{
		ID:           id,
		Name:         p.Name,
		Title:        p.Title,
		Summary:      p.Summary,
		Content:      p.Content,
		State:        p.State,
		Order:        p.Order,
		Template:     p.Template,
		Password:     p.Password,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Modified:     time.Now().Unix(),
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("putPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 更新文档内容
	if _, err := tx.Update(pp); err != nil {
		logs.Error("putPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 删除旧的关联内容
	sql := "DELETE FROM #relationships WHERE {postID}=?"
	if _, err := tx.Exec(true, sql, pp.ID); err != nil {
		logs.Error("putPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 添加新的关联
	rs := make([]*relationship, 0, len(p.Tags)+len(p.Cats))
	for _, v := range p.Tags {
		rs = append(rs, &relationship{MetaID: v, PostID: pp.ID})
	}
	for _, v := range p.Cats {
		rs = append(rs, &relationship{MetaID: v, PostID: pp.ID})
	}
	if err := tx.InsertMany(rs); err != nil {
		logs.Error("putPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logs.Error("putPost:", err)
		tx.Rollback()
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api delete /admin/api/posts/{id} 删除文章
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 no content
func adminDeletePost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("deletePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 删除文章
	sql := "DELETE FROM #posts WHERE {id}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("deletePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 删除评论
	sql = "DELETE FROM #comments WHERE {postID}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("deletePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	//删除关联数据
	sql = "DELETE FROM #relationships WHERE {postID}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("deletePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logs.Error("deletePost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api get /admin/api/posts 获取文章列表
// @apiQuery page int
// @apiQuery size int
// @apiQuery state int
// @apiGroup admin
//
// @apiSuccess ok 200
// @apiParam count int 符合条件的所有记录数量，不包含page和size条件
// @apiParam posts array 当前页的记录数量
func adminGetPosts(w http.ResponseWriter, r *http.Request) {
	var page, size, state int
	var ok bool
	if state, ok = core.QueryInt(w, r, "state", commentStateAll); !ok {
		return
	}

	sql := db.SQL().Table("#posts")
	if state != postStateAll {
		sql.And("{state}=?", state)
	}
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("adminGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	if size, ok = core.QueryInt(w, r, "size", opt.PageSize); !ok {
		return
	}
	sql.Limit(size, page*size)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("adminGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "posts": maps}, nil)
}

// @api get /admin/api/posts/{id} 获取某一篇文章的详细内容
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam id int id值
// @apiParam name string 唯一名称，可以为空
// @apiParam title string 标题
// @apiParam summary string 文章摘要
// @apiParam content string 文章内容
// @apiParam state int 状态
// @apiParam order int 排序
// @apiParam created int 创建时间
// @apiParam modified int 修改时间
// @apiParam template string 所使用的模板
// @apiParam password string 访问密码
// @apiParam allowPing bool 允许ping
// @apiParam allowComment bool 允许评论
// @apiParam tags array 关联的标签
// @apiParam cats array 关联的分类
func adminGetPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("adminGetPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	tags, err := getPostMetas(id, metaTypeTag)
	if err != nil {
		logs.Error("adminGetPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	cats, err := getPostMetas(id, metaTypeCat)
	if err != nil {
		logs.Error("adminGetPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	obj := &struct {
		ID           int64   `json:"id"`
		Name         string  `json:"name"`
		Title        string  `json:"title"`
		Summary      string  `json:"summary"`
		Content      string  `json:"content"`
		State        int     `json:"state"`
		Order        int     `json:"order"`
		Created      int64   `json:"created"`
		Modified     int64   `json:"modified"`
		Template     string  `json:"template"`
		Password     string  `json:"password"`
		AllowPing    bool    `json:"AllowPing"`
		AllowComment bool    `json:"AllowComment"`
		Tags         []int64 `json:"tags"`
		Cats         []int64 `json:"cats"`
	}{
		ID:           p.ID,
		Name:         p.Name,
		Title:        p.Title,
		Summary:      p.Summary,
		Content:      p.Content,
		State:        p.State,
		Order:        p.Order,
		Created:      p.Created,
		Modified:     p.Modified,
		Template:     p.Template,
		Password:     p.Password,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Tags:         tags,
		Cats:         cats,
	}
	core.RenderJSON(w, http.StatusOK, obj, nil)
}

// @api get /api/posts/{id} 获取某一文章的详细内容
// @apiGroup front
//
// @apiSuccess 200 ok
// @apiParam id int id值
// @apiParam name string 唯一名称，可以为空
// @apiParam title string 标题
// @apiParam content string 文章内容
// @apiParam created int 创建时间
// @apiParam modified int 修改时间
// @apiParam template string 所使用的模板
// @apiParam allowPing bool 允许ping
// @apiParam allowComment bool 允许评论
// @apiParam tags array 文章关联的标签
// @apiParam cats array 文章关联的类别
func frontGetPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("getPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != postStateNormal {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	tags, err := getPostMetas(id, metaTypeTag)
	if err != nil {
		logs.Error("getPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	cats, err := getPostMetas(id, metaTypeCat)
	if err != nil {
		logs.Error("getPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	obj := &struct {
		ID           int64   `json:"id"`
		Name         string  `json:"name"`
		Title        string  `json:"title"`
		Content      string  `json:"content"`
		Created      int64   `json:"created"`
		Modified     int64   `json:"modified"`
		Template     string  `json:"template"`
		AllowPing    bool    `json:"AllowPing"`
		AllowComment bool    `json:"AllowComment"`
		Tags         []int64 `json:"tags"`
		Cats         []int64 `json:"cats"`
	}{
		ID:           p.ID,
		Name:         p.Name,
		Title:        p.Title,
		Content:      p.Content,
		Created:      p.Created,
		Modified:     p.Modified,
		Template:     p.Template,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Tags:         tags,
		Cats:         cats,
	}
	core.RenderJSON(w, http.StatusOK, obj, nil)
}

// @api get /api/posts 获取前端可访问的文章列表
// @apiQuery page int
// @apiQuery size int
// @apiGroup front
//
// @apiSuccess 200 OK
func frontGetPosts(w http.ResponseWriter, r *http.Request) {
	sql := db.Where("{state}=?", postStateNormal).Table("#posts")
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("adminGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	var page, size int
	var ok bool
	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	if size, ok = core.QueryInt(w, r, "size", opt.PageSize); !ok {
		return
	}
	sql.Limit(size, page*size)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("adminGetPosts:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "posts": maps}, nil)
}

// @api get /api/posts/{id}/comments
// @apiQuery page int
// @apiQuery size int
// @apiGroup front
//
// @apiSuccess 200 OK
// @apiParam count int 当前文章的所有评论数量
// @apiParam comments array 当前页的评论
func frontGetPostComments(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("frontGetPostComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if p.State != postStateNormal {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	sql := db.Where("{postID}=?", id).
		And("{state}=?", commentStateApproved).
		Table("#comments")
	count, err := sql.Count(true)
	if err != nil {
		logs.Error("frontGetPostComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	var page, size int
	if page, ok = core.QueryInt(w, r, "page", 0); !ok {
		return
	}
	if size, ok = core.QueryInt(w, r, "size", opt.PageSize); !ok {
		return
	}
	sql.Limit(size, page*size)
	maps, err := sql.SelectMap(true, "*")
	if err != nil {
		logs.Error("getComments:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"count": count, "comments": maps}, nil)
}

// @api post /api/posts/{id}/comments 提交新评论
// @apiGroup front
//
// @apiRequest json
// @apiParam parent int 评论的父级内容
// @apiParam postID int 评论的文章
// @apiParam content string 评论的内容
// @apiParam authorName string 评论的作者
// @apiParam authorURL string 评论作者的网站地址，可为空
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

	p := &post{ID: c.PostID}
	if err := db.Select(p); err != nil {
		logs.Error("postPostComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if (len(p.Title) == 0 && len(p.Content) == 0) || p.State != postStateNormal {
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
		logs.Error("postComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	c.AuthorURL = u.Scheme + ":" + u.Host

	c.Content = html.EscapeString(c.Content)
	c.Content = strings.Replace(c.Content, "\n", "<br />", -1)

	comm := &comment{
		PostID:      c.PostID,
		Parent:      c.Parent,
		AuthorURL:   c.AuthorURL,
		AuthorName:  c.AuthorName,
		AuthorEmail: c.AuthorEmail,
		Content:     c.Content,
		Created:     time.Now().Unix(),
		State:       commentStateWaiting,
		IP:          r.RemoteAddr,
		Agent:       r.UserAgent(),
		IsAdmin:     false,
	}
	if _, err := db.Insert(comm); err != nil {
		logs.Error("postComment:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, nil, nil)
}

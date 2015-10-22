// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api post /admin/api/posts/{id}/published 将一篇文章的状态改为已发布
// @apiParam id int 文章的id
// @apiGroup admin
//
// @apiSuccess 201 Created
func adminSetPostPublished(w http.ResponseWriter, r *http.Request) {
	adminSetPostState(w, r, models.PostStatePublished)
}

// @api post /admin/api/posts/{id}/draft 将一篇文章的状态改为草稿
// @apiParam id int 文章的id
// @apiGroup admin
//
// @apiSuccess 201 Created
func adminSetPostDraft(w http.ResponseWriter, r *http.Request) {
	adminSetPostState(w, r, models.PostStateDraft)
}

func adminSetPostState(w http.ResponseWriter, r *http.Request, state int) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("adminSetPostState:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	// 不可能存在状态值为0的文章，出现此值，表明数据库没有该条记录
	if p.State == models.PostStateAll {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	p = &models.Post{ID: id, State: state}
	if _, err := db.Update(p); err != nil {
		logs.Error("adminSetPostState:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, "{}", nil)
}

// @api get /admin/api/posts/count 获取各种状态下的文章数量
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam all       int 评论的总量
// @apiParam draft     int 等待审核的评论数量
// @apiParam published int 垃圾评论的数量
func adminGetPostsCount(w http.ResponseWriter, r *http.Request) {
	sql := "SELECT {state}, count(*) AS cnt FROM #posts GROUP BY {state}"
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("adminGetPostsCount:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	defer rows.Close()

	maps, err := fetch.MapString(false, rows)
	if err != nil {
		logs.Error("adminGetPostsCount:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := map[string]int{
		"all":       0,
		"draft":     0,
		"published": 0,
	}
	count := 0
	for _, v := range maps {
		num, err := strconv.Atoi(v["cnt"])
		if err != nil {
			logs.Error("adminGetCommentsCount:", err)
			core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
			return
		}
		count += num
		switch v["state"] {
		case "1":
			data["published"] = num
		case "2":
			data["draft"] = num
		default:
			logs.Error("adminGetPostsCount: 未知的文章状态:", v["state"])
		}
	}
	data["all"] = count // 所有评论的数量
	core.RenderJSON(w, http.StatusOK, data, nil)
}

// @api post /admin/api/posts 新建文章
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name         string 唯一名称，可以为空
// @apiParam title        string 标题
// @apiParam summary      string 文章摘要
// @apiParam content      string 文章内容
// @apiParam state        int    状态
// @apiParam order        int    排序
// @apiParam template     string 所使用的模板
// @apiParam allowPing    bool   允许ping
// @apiParam allowComment bool   允许评论
// @apiParam tags         string 关联的标签，多个标签名称以逗号分隔
//
// @apiSuccess 201 created
func adminPostPost(w http.ResponseWriter, r *http.Request) {
	p := &struct {
		Name         string `json:"name"`
		Title        string `json:"title"`
		Summary      string `json:"summary"`
		Content      string `json:"content"`
		State        int    `json:"state"`
		Order        int    `json:"order"`
		Template     string `json:"template"`
		AllowPing    bool   `json:"allowPing"`
		AllowComment bool   `json:"allowComment"`
		Tags         string `json:"tags"`
	}{}

	if !core.ReadJSON(w, r, p) {
		return
	}

	t := time.Now().Unix()
	pp := &models.Post{
		Name:         p.Name,
		Title:        p.Title,
		Summary:      p.Summary,
		Content:      p.Content,
		State:        p.State,
		Order:        p.Order,
		Template:     p.Template,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Created:      t,
		Modified:     t,
	}

	tags, err := getTagsID(p.Tags)
	if err != nil {
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 插入文章
	result, err := tx.Insert(pp)
	if err != nil {
		tx.Rollback()
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 插入relationship
	rs := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		rs = append(rs, &models.Relationship{PostID: postID, TagID: tag})
	}
	if err := tx.MultInsert(rs...); err != nil {
		tx.Rollback()
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// commit
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logs.Error("adminPostPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, "{}", nil)
}

// @api put /admin/api/posts/{id} 修改文章
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name         string 唯一名称，可以为空
// @apiParam title        string 标题
// @apiParam summary      string 文章摘要
// @apiParam content      string 文章内容
// @apiParam state        int    状态
// @apiParam order        int    排序
// @apiParam template     string 所使用的模板
// @apiParam allowPing    bool   允许ping
// @apiParam allowComment bool   允许评论
// @apiParam tags         string 关联的标签，多个标签名称以逗号分隔
//
// @apiSuccess 200 no content
func adminPutPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &struct {
		Name         string `json:"name"`
		Title        string `json:"title"`
		Summary      string `json:"summary"`
		Content      string `json:"content"`
		State        int    `json:"state"`
		Order        int    `json:"order"`
		Template     string `json:"template"`
		AllowPing    bool   `json:"allowPing"`
		AllowComment bool   `json:"allowComment"`
		Tags         string `json:"tags"`
	}{}
	if !core.ReadJSON(w, r, p) {
		return
	}

	pp := &models.Post{
		ID:           id,
		Name:         p.Name,
		Title:        p.Title,
		Summary:      p.Summary,
		Content:      p.Content,
		State:        p.State,
		Order:        p.Order,
		Template:     p.Template,
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Modified:     time.Now().Unix(),
	}
	tags, err := getTagsID(p.Tags)
	if err != nil {
		logs.Error("adminPostPost-0:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("adminPutPost-1:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 更新文档内容
	if _, err := tx.Update(pp); err != nil {
		logs.Error("adminPutPost-2:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 删除旧的关联内容
	sql := "DELETE FROM #relationships WHERE {postID}=?"
	if _, err := tx.Exec(true, sql, pp.ID); err != nil {
		logs.Error("adminPutPost-3:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		tx.Rollback()
		return
	}

	// 添加新的关联
	if len(p.Tags) > 0 {
		rs := make([]interface{}, 0, len(p.Tags))
		for _, tag := range tags {
			rs = append(rs, &models.Relationship{TagID: tag, PostID: pp.ID})
		}
		if err := tx.MultInsert(rs...); err != nil {
			logs.Error("adminPutPost-4:", err)
			core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
			tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logs.Error("adminPutPost-5:", err)
		tx.Rollback()
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 将一串标签名转换成id
// names为一种由标签名组成的字符串，名称之间由逗号分隔。
func getTagsID(names string) ([]int64, error) {
	name := strings.Split(names, ",")
	if len(name) == 0 {
		return nil, nil
	}

	cond := strings.Repeat("?,", len(name))
	sql := "SELECT {id} FROM #tags WHERE {title} IN(" + cond[:len(cond)-1] + ")"
	args := make([]interface{}, 0, len(name))
	for _, v := range name {
		args = append(args, v)
	}
	rows, err := db.Query(true, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "id", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, str := range maps {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}

	return ret, nil
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
// @apiQuery page  int 页码，从0开始
// @apiQuery size  int 显示尺寸
// @apiQuery state int 状态
// @apiGroup admin
//
// @apiSuccess ok 200
// @apiParam count int   符合条件的所有记录数量，不包含page和size条件
// @apiParam posts array 当前页的记录数量
func adminGetPosts(w http.ResponseWriter, r *http.Request) {
	var page, size, state int
	var ok bool
	if state, ok = core.QueryInt(w, r, "state", models.CommentStateAll); !ok {
		return
	}

	sql := db.SQL().Table("#posts")
	if state != models.PostStateAll {
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
	maps, err := sql.SelectMapString(true, "*")
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
// @apiParam id           int    id值
// @apiParam name         string 唯一名称，可以为空
// @apiParam title        string 标题
// @apiParam summary      string 文章摘要
// @apiParam content      string 文章内容
// @apiParam state        int    状态
// @apiParam order        int    排序
// @apiParam created      int    创建时间
// @apiParam modified     int    修改时间
// @apiParam template     string 所使用的模板
// @apiParam allowPing    bool   允许ping
// @apiParam allowComment bool   允许评论
// @apiParam tags         string 关联的标签，多个标签以逗号分隔。
func adminGetPost(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	p := &models.Post{ID: id}
	if err := db.Select(p); err != nil {
		logs.Error("adminGetPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	tags, err := getPostTagsName(id)
	if err != nil {
		logs.Error("adminGetPost:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	obj := &struct {
		ID           int64  `json:"id"`
		Name         string `json:"name"`
		Title        string `json:"title"`
		Summary      string `json:"summary"`
		Content      string `json:"content"`
		State        int    `json:"state"`
		Order        int    `json:"order"`
		Created      int64  `json:"created"`
		Modified     int64  `json:"modified"`
		Template     string `json:"template"`
		AllowPing    bool   `json:"allowPing"`
		AllowComment bool   `json:"allowComment"`
		Tags         string `json:"tags"`
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
		AllowPing:    p.AllowPing,
		AllowComment: p.AllowComment,
		Tags:         tags,
	}
	core.RenderJSON(w, http.StatusOK, obj, nil)
}

// 获取与某post相关联的标签
func getPostTagsName(postID int64) (string, error) {
	sql := `SELECT t.{title} FROM #relationships AS rs
	LEFT JOIN #tags AS t ON t.{id}=rs.{tagID}
	WHERE rs.{postID}=?`
	rows, err := db.Query(true, sql, postID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "title", rows)
	if err != nil {
		return "", err
	}

	return strings.Join(maps, ","), nil
}

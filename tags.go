// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api put /admin/api/tags/{id}/merge 将指定的标签合并到当前标签
// @apiGroup admin
//
// @apiRequest json
// @apiParam tags array 所有需要合并的标签ID列表。
// @apiExample json
// {"tags": [1,2,3] }
func adminPutTagMerge(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// @api get /admin/api/tags 获取所有标签信息
// @apiGroup admin
//
// @apiRequest json
// @apiheader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam tags array 所有分类的列表
func adminGetTags(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT m.{name},m.{title},m.{description},m.{id},count(r.{tagID}) AS {count}
			FROM #tags AS m
			LEFT JOIN #relationships AS r ON m.{id}=r.{tagID}
			GROUP BY m.{id}`
	rows, err := db.Query(true, sql)
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"tags": maps}, nil)
}

// @api put /admin/api/tags/{id} 修改某id的标签内容
// @apiParam id int 需要修改的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 204 no content
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPutTag(w http.ResponseWriter, r *http.Request) {
	m := &models.Tag{}
	if !core.ReadJSON(w, r, m) {
		return
	}

	var ok bool
	m.ID, ok = core.ParamID(w, r, "id")
	if !ok {
		return
	}

	exists, err := tagNameIsExists(m, typ)
	if err != nil {
		logs.Error("adminPutTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	errs := &core.ErrorResult{Message: "格式错误"}
	if exists {
		errs.Detail["name"] = "已有同名字体段"
	}

	if len(m.Title) == 0 {
		errs.Detail["title"] = "标题不能为空"
	}

	if len(errs.Detail) > 0 {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	if _, err := db.Update(m); err != nil {
		logs.Error("adminPutTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api post /admin/api/tags 添加新标签
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 201 created
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPostTag(w http.ResponseWriter, r *http.Request) {
	m := &models.Tag{}
	if !core.ReadJSON(w, r, m) {
		return
	}

	errs := &core.ErrorResult{Message: "格式错误"}
	if m.ID > 0 {
		errs.Detail["id"] = "不允许的字段"
	}
	m.ID = 0

	exists, err := tagNameIsExists(m, typ)
	if err != nil {
		logs.Error("adminPostTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if exists {
		errs.Detail["name"] = "已有同名字体段"
	}

	if len(m.Title) == 0 {
		errs.Detail["title"] = "标题不能为空"
	}

	if len(errs.Detail) > 0 {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}
	m.Type = typ

	if _, err := db.Insert(m); err != nil {
		logs.Error("adminPostTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, "{}", nil)
}

// @api delete /admin/api/tags/{id} 删除该id的标签
// @apiParam id int 需要删除的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 no content
func adminDeleteTag(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if _, err := tx.Delete(&models.Tag{ID: id}); err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 删除与之对应的关联数据。
	sql := "DELETE FROM #relationships WHERE {tagID}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("adminDeleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := tx.Commit(); err != nil {
		logs.Error("adminDeleteMeta:", err)
		tx.Rollback()
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 是否存在相同name的tag
func tagNameIsExists(m *models.Meta, typ int) (bool, error) {
	sql := db.Where("{name}=?", m.Name).And("{type}=?", typ).Table("#tags")
	maps, err := sql.SelectMapString(true, "id")
	if err != nil {
		return false, err
	}

	if len(maps) == 0 {
		return false, nil
	}
	if len(maps) > 1 {
		return true, nil
	}

	id, err := strconv.ParseInt(maps[0]["id"], 10, 64)
	println(maps[0]["id"])
	if err != nil {
		return false, err
	}
	return id != m.ID, nil
}

// @api get /admin/api/tags/{id} 获取指定id的标签内容
// @apiParam id int 标签的id
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam id int 标签的id
// @apiParam name string 标签的唯一名称，可能为空
// @apiParam title string 标签名称
// @apiParam description string 对标签的详细描述
func adminGetTag(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	m := &models.Meta{ID: id}
	if err := db.Select(m); err != nil {
		logs.Error("adminGetTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := &struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}{
		ID:          m.ID,
		Name:        m.Name,
		Title:       m.Title,
		Description: m.Description,
	}
	core.RenderJSON(w, http.StatusOK, data, nil)
}

// 获取与某post相关联的标签或是分类
func getPostMetas(postID int64, mtype int) ([]int64, error) {
	sql := `SELECT rs.{tagID} FROM #relationships AS rs
	LEFT JOIN #tags AS m ON m.{id}=rs.{tagID}
	WHERE rs.{postID}=? AND m.{type}=?`
	rows, err := db.Query(true, sql, postID, mtype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "tagID", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, v := range maps {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}
	return ret, nil
}

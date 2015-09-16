// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
)

// @api get /admin/api/themes 获取所有主题列表
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam themes array 所有主题列表
func adminGetThemes(w http.ResponseWriter, r *http.Request) {
	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"themes": themes.Themes()}, nil)
}

// @api patch /admin/api/themes/current 更改当前的主题
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam value string 新值
//
// @apiSuccess 200 OK
func adminPatchTheme(w http.ResponseWriter, r *http.Request) {
	o := &models.Option{Key: "theme"}
	if !core.ReadJSON(w, r, o) {
		return
	}

	if o.Key != "theme" || len(o.Group) > 0 { // 提交了额外的数据内容
		core.RenderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	if err := patchOption(o); err != nil {
		logs.Error("adminPatchTheme:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := themes.LoadTheme(o.Value); err != nil {
		logs.Error("adminPatchTheme:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/themes"
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

// @api get /admin/api/themes/current 获取当前的主题信息。
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam theme string 主题名称
func adminGetCurrentTheme(w http.ResponseWriter, r *http.Request) {
	core.RenderJSON(w, http.StatusOK, map[string]string{"theme": opt.Theme}, nil)
}

// @api put /admin/api/themes/current 更改当前的主题
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam value string 新值
//
// @apiSuccess 200 OK
func adminPutCurrentTheme(w http.ResponseWriter, r *http.Request) {
	v := &struct {
		Value string `json:"value"`
	}{}
	if !core.ReadJSON(w, r, v) {
		return
	}

	if len(v.Value) == 0 {
		core.RenderJSON(w, http.StatusBadRequest, &core.ErrorResult{Message: "必须指定一个值！"}, nil)
		return
	}

	o := &models.Option{Key: "theme", Value: v.Value}
	if err := patchOption(o); err != nil {
		logs.Error("adminPutTheme:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := themes.Switch(o.Value); err != nil {
		logs.Error("adminPutTheme:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

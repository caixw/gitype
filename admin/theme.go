// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"github.com/caixw/typing/themes"
	"github.com/caixw/typing/util"
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
	util.RenderJSON(w, http.StatusOK, map[string]interface{}{"themes": themes.Themes()}, nil)
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
	util.RenderJSON(w, http.StatusOK, map[string]string{"theme": opt.Theme}, nil)
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
	if !util.ReadJSON(w, r, v) {
		return
	}

	if len(v.Value) == 0 {
		util.RenderJSON(w, http.StatusBadRequest, &util.ErrorResult{Message: "必须指定一个值！"}, nil)
		return
	}

	if err := opt.Set(db, "theme", v.Value, false); err != nil {
		logs.Error("adminPutTheme:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := themes.Switch(v.Value); err != nil {
		logs.Error("adminPutTheme:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	lastUpdated()
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

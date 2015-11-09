// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/caixw/typing/core"
)

// @api get /admin/api/media 获取所有的文件列表
// @apiQuery parent string 上一级目录
// @apiGroup admin
//
// @apiSuccess 200 成功获取列表
// @apiParam files array 文件列表
func adminGetMedia(w http.ResponseWriter, r *http.Request) {
	parent, ok := core.ParamString(w, r, "parent")
	if !ok {
		return
	}
	if strings.Index(parent, "..") >= 0 {
		core.RenderJSON(w, http.StatusBadRequest, &core.ErrorResult{Message: "格式错误"}, nil)
		return
	}
	parent = core.Cfg.UploadDir + parent

	list, err := ioutil.ReadDir(parent)
	if err != nil {
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"files": list}, nil)
}

// @api post /admin/api/media 上传媒体文件
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam media file 文件内容
//
// @apiSuccess 201 文件上传成功
func adminPostMedia(w http.ResponseWriter, r *http.Request) {
	files, err := u.Do("media", r)
	if err != nil {
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusCreated, map[string]interface{}{"media": files[0]}, nil)
}

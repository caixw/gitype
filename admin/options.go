// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/sitemap"
	"github.com/issue9/logs"
)

// @api patch /admin/api/options/{key} 修改设置项的值
// @apiParam key string 需要修改项的key
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam value string 新值
// @apiExample json
// { "value": "abcdef" }
// @apiSuccess 204 no content
func adminPatchOption(w http.ResponseWriter, r *http.Request) {
	key, ok := core.ParamString(w, r, "key")
	if !ok {
		return
	}

	o := &models.Option{Key: key}
	cnt, err := db.Count(o)
	if err != nil {
		logs.Error("patchOption:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if cnt == 0 {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	if !core.ReadJSON(w, r, o) {
		return
	}

	if o.Key != key || len(o.Group) > 0 { // 提交了额外的数据内容
		core.RenderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	if err := patchOption(o); err != nil {
		logs.Error("patchOption:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

func patchOption(o *models.Option) error {
	// 更新数据库中的值
	if _, err := db.Update(o); err != nil {
		return err
	}

	// 更新opt变量中的值
	return opt.UpdateFromOption(o)
}

// @api get /admin/api/options/{key} 更新设置项的值，不能更新password字段。
// @apiParam key string 名称
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 ok
// @api value any 设置项的值
// @apiExample json
// { "value": "20" }
func adminGetOption(w http.ResponseWriter, r *http.Request) {
	key, ok := core.ParamString(w, r, "key")
	if !ok {
		return
	}

	if key == "password" {
		core.RenderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	val, found := opt.GetValueByKey(key)
	if !found {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"value": val}, nil)
}

// @api put /admin/api/sitemap 重新生成sitemap
// @apiGroup admin
//
// @apiSuccess 200 Ok
func adminPutSitemap(w http.ResponseWriter, r *http.Request) {
	err := sitemap.Build(db, opt)
	if err != nil {
		logs.Error(err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, "{}", nil)
}

// @api get /admin/api/state 获取当前网站的些基本状态
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam posts            int 文章的数量
// @apiParam draftPosts       int 草稿的数量
// @apiParam normalPosts      int 正式文章的数量
// @apiParam comments         int 评论数量
// @apiParam waitingComments  int 待审评论数量
// @apiParam spamComments     int 垃圾评论数量
// @apiParam approvedComments int 已审评论数量
// @apiParam lastLogin        int 最后次登录时间
// @apiParam lastPost         int 最后次发表文章的时间
// @apiParam lastIP           string 最后次登录的IP
// @apiParam lastAgent        string 最后次登录的浏览器相关资料
// @apiParam screenName       string 用户的当前昵称
func adminGetState(w http.ResponseWriter, r *http.Request) {
	// TODO
}

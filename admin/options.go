// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/util"
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
	key, ok := util.ParamString(w, r, "key")
	if !ok {
		return
	}

	if _, found := opt.Get(key); !found {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	data := &struct {
		Value string `json:"value"`
	}{}
	if !util.ReadJSON(w, r, data) {
		return
	}

	if err := opt.Set(db, key, data.Value, false); err != nil {
		logs.Error("adminPatchOption:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	lastUpdated()
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api get /admin/api/options/{key} 获取设置项的值，不能获取password字段。
// @apiParam key string 名称
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 ok
// @api value any 设置项的值
// @apiExample json
// { "value": "20" }
func adminGetOption(w http.ResponseWriter, r *http.Request) {
	key, ok := util.ParamString(w, r, "key")
	if !ok {
		return
	}

	if key == "password" {
		util.RenderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	val, found := opt.Get(key)
	if !found {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, map[string]interface{}{"value": val}, nil)
}

// @api put /admin/api/sitemap 重新生成sitemap
// @apiGroup admin
//
// @apiSuccess 200 Ok
func adminPutSitemap(w http.ResponseWriter, r *http.Request) {
	err := feed.BuildSitemap()
	if err != nil {
		logs.Error(err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	lastUpdated()
	util.RenderJSON(w, http.StatusOK, "{}", nil)
}

// @api get /admin/api/stat 获取当前网站的些基本状态
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam posts            int 文章的数量
// @apiParam draftPosts       int 草稿的数量
// @apiParam publishedPosts   int 正式文章的数量
// @apiParam comments         int 评论数量
// @apiParam waitingComments  int 待审评论数量
// @apiParam spamComments     int 垃圾评论数量
// @apiParam approvedComments int 已审评论数量
// @apiParam lastLogin        int 最后次登录时间
// @apiParam lastPost         int 最后次发表文章的时间
// @apiParam lastIP           string 最后次登录的IP
// @apiParam lastAgent        string 最后次登录的浏览器相关资料
// @apiParam screenName       string 用户的当前昵称
func adminGetStat(w http.ResponseWriter, r *http.Request) {
	// TODO
}

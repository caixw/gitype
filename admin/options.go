// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/util"
	"github.com/issue9/logs"
	"github.com/issue9/web"
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
//
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
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	val, found := opt.Get(key)
	if !found {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, map[string]interface{}{"value": val}, nil)
}

// @api put /admin/api/feed/sitemap 重新生成sitemap
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 Ok
func adminPutSitemap(w http.ResponseWriter, r *http.Request) {
	err := feed.BuildSitemap()
	if err != nil {
		logs.Error(err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, "{}", nil)
}

// @api put /admin/api/feed/rss 重新生成rss
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 Ok
func adminPutRss(w http.ResponseWriter, r *http.Request) {
	err := feed.BuildRss()
	if err != nil {
		logs.Error(err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, "{}", nil)
}

// @api put /admin/api/feed/atom 重新生成atom
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 Ok
func adminPutAtom(w http.ResponseWriter, r *http.Request) {
	err := feed.BuildAtom()
	if err != nil {
		logs.Error(err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	util.RenderJSON(w, http.StatusOK, "{}", nil)
}

// @api get /admin/api/state 获取当前网站的些基本状态
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam posts            int 文章的数量
// @apiParam draftPosts       int 草稿的数量
// @apiParam publishedPosts   int 正式文章的数量
// @apiParam comments         int 评论数量
// @apiParam waitingComments  int 待审评论数量
// @apiParam spamComments     int 垃圾评论数量
// @apiParam approvedComments int 已审评论数量
// @apiParam screenName       string 用户的当前昵称
func adminGetState(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"posts":            stat.PostsSize,
		"draftPosts":       stat.DraftPostsSize,
		"publishedPosts":   stat.PublishedPostsSize,
		"comments":         stat.CommentsSize,
		"waitingComments":  stat.WaitingCommentsSize,
		"spamComments":     stat.SpamCommentsSize,
		"approvedComments": stat.ApprovedCommentsSize,
		"screenName":       opt.ScreenName,
	}

	util.RenderJSON(w, http.StatusOK, data, nil)
}

// @api get /admin/api/modules 获取的模块列表
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam modules array 模块列表的数组
func adminGetModules(w http.ResponseWriter, r *http.Request) {
	modules := web.Modules()
	data := make([]map[string]interface{}, 0, len(modules))
	for _, v := range modules {
		data = append(data, map[string]interface{}{
			"name":      v.Name,
			"isRunning": v.IsRunning(),
		})
	}
	util.RenderJSON(w, http.StatusOK, map[string]interface{}{"modules": data}, nil)
}

// @api put /admin/api/modules/{name}/start 启动一个模块
// @apiParam name string 模块名称
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 OK
func adminPutModuleStart(w http.ResponseWriter, r *http.Request) {
	m := getModule(w, r)
	if m == nil {
		return
	}

	m.Start()
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api put /admin/api/modules/{name}/stop 停止一个模块
// @apiParam name string 模块名称
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 OK
func adminPutModuleStop(w http.ResponseWriter, r *http.Request) {
	m := getModule(w, r)
	if m == nil {
		return
	}

	m.Stop()
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

func getModule(w http.ResponseWriter, r *http.Request) *web.Module {
	name, ok := util.ParamString(w, r, "name")

	if !ok {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return nil
	}

	m := web.GetModule(name)
	if m == nil {
		util.RenderJSON(w, http.StatusNotFound, nil, nil)
		return nil
	}

	return m
}

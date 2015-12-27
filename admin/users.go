// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"crypto/rand"
	"io"
	"net/http"

	"github.com/caixw/typing/models"
	"github.com/caixw/typing/util"
	"github.com/issue9/logs"
)

// 记录登录后的token值
var token string

// @api post /admin/api/login 登录
// @apiGroup admin
//
// @apiRequest json
// @apiParam password string 登录密码
// @apiExample json
// { "password": "12345" }
//
// @apiSuccess 201
// @apiHeader Cache-Control:no-cache
// @apiHeader Pragma:no-cache
// @apiParam token string 登录凭证；
// @apiExample json
// { "token":  "adfwerqeqaeqe313aa" }
func adminPostLogin(w http.ResponseWriter, r *http.Request) {
	inst := &struct {
		Password string `json:"password"`
	}{}
	if !util.ReadJSON(w, r, inst) {
		return
	}

	if cfg.Password(inst.Password) != opt.Password {
		util.RenderJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	ret := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, ret)
	if err != nil {
		logs.Error("login:无法产生一个随机的token", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if n == 0 {
		logs.Error("login:无法产生一个随机的token")
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	token = util.MD5(string(ret))
	util.RenderJSON(w, http.StatusCreated, map[string]string{"token": token}, nil)
}

// @api delete /admin/api/login 注销当前用户的登录
// @apiGroup admin
// @apiRequest json
// @apiHeader Authorization xxxxx
// @apiSuccess 204 no content
func adminDeleteLogin(w http.ResponseWriter, r *http.Request) {
	token = ""
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api put /admin/api/password 理发密码
// @apiGroup admin
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam old string 旧密码
// @apiParam new string 新密码
// @apiExample json
// {
//     "old": "123",
//     "new": "456"
// }
//
// @apiSuccess 204 no content
func adminChangePassword(w http.ResponseWriter, r *http.Request) {
	l := &struct {
		Old string `json:"old"`
		New string `json:"new"`
	}{}

	if !util.ReadJSON(w, r, l) {
		return
	}

	errs := &util.ErrorResult{Message: "提交数据错误", Detail: map[string]string{}}
	if len(l.New) == 0 {
		errs.Add("new", "新密码不能为空")
	}
	if opt.Password != cfg.Password(l.Old) {
		errs.Add("old", "旧密码错误")
	}
	if len(errs.Detail) > 0 {
		util.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	o := &models.Option{Key: "password", Value: cfg.Password(l.New)}
	if _, err := db.Update(o); err != nil {
		logs.Error("adminChangePassword:", err)
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	opt.Password = o.Value
	util.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 验证后台登录信息
func loginHandlerFunc(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != token {
			util.RenderJSON(w, http.StatusUnauthorized, nil, nil)
			return
		}
		f(w, r)
	}
	return http.HandlerFunc(h)
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package admin

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/caixw/typing/util"
	"github.com/issue9/logs"
)

// 最大的登录日志数量
const maxLoginLogs = 5

type login struct {
	IP    string `json:"ip,omitempty"`
	Agent string `json:"agent,omitempty"`
	Time  string `json:"time,omitempty"`
}

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

	if app.Password(inst.Password) != opt.Password {
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
	if len(token) == 0 {
		logs.Error("login:无法正确生成登录的token")
		util.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 记录日志出错，仅输出错误内容，但不返回500错误。
	if err = writeLastLogs(r); err != nil {
		logs.Error("login:", err)
	}

	util.RenderJSON(w, http.StatusCreated, map[string]string{"token": token}, nil)
}

func writeLastLogs(r *http.Request) error {
	ls := make([]*login, 0, maxLoginLogs)
	if err := json.Unmarshal([]byte(opt.Last), &ls); err != nil {
		return err
	}

	l := &login{
		IP:    r.RemoteAddr,
		Agent: r.UserAgent(),
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}

	if len(ls) >= maxLoginLogs { // 去掉最后一条记录
		ls = ls[:maxLoginLogs-1]
	}
	//lss := make([]*login, 0, maxLoginLogs)
	//lss = append(lss, l)
	lss := []*login{l}
	lss = append(lss, ls...)
	bs, err := json.Marshal(lss)
	if err != nil {
		return err
	}

	if err := app.SetOption(db, "last", string(bs), true); err != nil {
		return err
	}
	logs.Infof("登录信息：IP:%v;Agent:%v;Time:%v\n", l.IP, l.Agent, l.Time)

	return nil
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
	if opt.Password != app.Password(l.Old) {
		errs.Add("old", "旧密码错误")
	}
	if len(errs.Detail) > 0 {
		util.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	o := &models.Option{Key: "password", Value: app.Password(l.New)}
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

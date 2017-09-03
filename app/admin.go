// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/caixw/typing/app/admin"
	"github.com/issue9/logs"
)

// 初始化控制台相关内容
func (a *app) initAdmin() (err error) {
	a.adminTpl, err = template.New("admin").Parse(admin.AdminHTML)
	if err != nil {
		return
	}

	a.mux.GetFunc(a.conf.AdminURL, a.getAdminPage).
		PostFunc(a.conf.AdminURL, a.postAdminPage)
	return nil
}

func (a *app) postAdminPage(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("password") != a.conf.AdminPassword {
		a.renderAdminPage(w, r, "密码错误！")
		return
	}

	if err := a.pull(); err.status >= 400 {
		logs.Error(err.message)
		statusError(w, err.status)
		return
	}

	a.renderAdminPage(w, r, "")
}

// 一个简单的后台页面，可用来手动更新数据。
func (a *app) getAdminPage(w http.ResponseWriter, r *http.Request) {
	a.renderAdminPage(w, r, "")
}

// message 表示出错的信息，空值表示没有错误
func (a *app) renderAdminPage(w http.ResponseWriter, r *http.Request, message string) {
	home := strings.TrimSuffix(r.URL.Path, a.conf.AdminURL)
	if len(home) == 0 {
		home = "/"
	}

	s := map[string]interface{}{
		"lastUpdate": a.client.Created.Format(time.RFC3339),
		"homeURL":    home,
		"message":    message,
	}

	if err := a.adminTpl.Execute(w, s); err != nil {
		logs.Error(err)
		statusError(w, http.StatusInternalServerError)
	}
}

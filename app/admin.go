// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"html/template"
	"net/http"
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
	if r.FormValue("password") == a.conf.AdminPassword {
		if err := a.reload(); err != nil {
			logs.Error(err)
			statusError(w, http.StatusInternalServerError)
			return
		}
	}

	a.getAdminPage(w, r)
}

// 一个简单的后台页面，可用来手动更新数据。
func (a *app) getAdminPage(w http.ResponseWriter, r *http.Request) {
	s := map[string]interface{}{
		"lastUpdate": time.Unix(a.buf.Updated, 0).Format(a.buf.Data.Config.LongDateFormat),
		"homeURL":    a.buf.Data.Config.URL,
	}

	if err := a.adminTpl.Execute(w, s); err != nil {
		logs.Error(err)
		statusError(w, http.StatusInternalServerError)
	}
}

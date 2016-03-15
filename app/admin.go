// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"html/template"
	"net/http"
	"os/exec"
	"time"

	"github.com/issue9/logs"
	"github.com/issue9/utils"
	"github.com/issue9/web"
)

// 初始化控制台相关内容
func (a *app) initAdmin() (err error) {
	a.adminTpl, err = template.New("admin").Parse(adminHTML)
	if err != nil {
		return
	}

	admin, err := web.NewModule("admin")
	if err != nil {
		return err
	}

	admin.GetFunc(a.conf.AdminURL, a.getAdminPage).
		PostFunc(a.conf.AdminURL, a.postAdminPage).
		PostFunc(a.conf.WebhooksURL, a.postWebhooks)
	return nil
}

func (a *app) postWebhooks(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Unix()

	if now-a.conf.WebhooksUpdateFreq < a.updated { // 时间太短，不接受更新
		logs.Info("更新过于频繁，被中止！")
		return
	}

	cmd := &exec.Cmd{
		Path: "git",
		Dir:  a.path.Data,
	}

	if utils.FileExists(cmd.Dir) {
		cmd.Args = []string{"pull"}
	} else {
		cmd.Args = []string{"clone", a.conf.RepoURL}
	}

	if err := cmd.Run(); err != nil {
		logs.Error("a.postWebhooks:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *app) postAdminPage(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("password") == a.conf.AdminPassword {
		if err := a.reload(); err != nil {
			logs.Error("app.postAdminPage:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	a.getAdminPage(w, r)
}

func (a *app) getAdminPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"lastUpdate": time.Unix(a.updated, 0).Format("2006-01-02 15:04:05+0700"),
	}

	if err := a.adminTpl.Execute(w, data); err != nil {
		logs.Error("app.getAdminPage:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

const adminHTML = `<!DOCTYPE html>
<html>
	<meta charset="utf-8" />
	<title>typing控制面板</title>
	<body style="text-align:center">
	<div style="margin:auto;text-align:left;width:30em">
		<h1>控制面板</h1>
		<p>
			<span>最后更新时间:</span>{{.lastUpdate}}
		</p>
		<form action="" method="POST">
			<p>
				<input type="password" name="password" placeholder="密码" />
				<button type="submit">重新生成</button>
			</p>
		</form>
	</div>
	</body>
</html>`

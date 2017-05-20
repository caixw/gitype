// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/caixw/typing/app/static"
	"github.com/issue9/logs"
	"github.com/issue9/utils"
)

// 初始化控制台相关内容
func (a *app) initAdmin() (err error) {
	a.adminTpl, err = template.New("admin").Parse(static.AdminHTML)
	if err != nil {
		return
	}

	a.mux.GetFunc(a.conf.AdminURL, a.getAdminPage).
		PostFunc(a.conf.AdminURL, a.postAdminPage).
		PostFunc(a.conf.WebhooksURL, a.postWebhooks)
	return nil
}

// 将一个log.Logger封装成io.Writer
type logW struct {
	l *log.Logger
}

func (w *logW) Write(bs []byte) (int, error) {
	w.l.Print(string(bs))
	return len(bs), nil
}

// 通过webhooks来更新内容
func (a *app) postWebhooks(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Unix()

	if now-a.conf.WebhooksUpdateFreq < a.updated { // 时间太短，不接受更新
		logs.Info("更新过于频繁，被中止！")
		return
	}

	var cmd *exec.Cmd
	if utils.FileExists(a.path.Data) {
		cmd = exec.Command("git", "pull")
		cmd.Dir = a.path.Data
	} else {
		cmd = exec.Command("git", "clone", a.conf.RepoURL, a.path.Data)
		cmd.Dir = a.path.Root
	}

	cmd.Stderr = &logW{l: logs.ERROR()}
	cmd.Stdout = &logW{l: logs.INFO()}
	if err := cmd.Run(); err != nil {
		logs.Error("a.postWebhooks:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.reload(); err != nil {
		logs.Error("app.postWebhooks:", err)
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

// 一个简单的后台页面，可用来手动更新加载新数据。
//
// 若数据不是通过github来管理的，可通过此方法来手动更新数据。
func (a *app) getAdminPage(w http.ResponseWriter, r *http.Request) {
	var homeURL string

	// data 有可能加载失败
	if a.data != nil {
		homeURL = a.data.Config.URL
	}
	data := map[string]interface{}{
		"lastUpdate": time.Unix(a.updated, 0).Format("2006-01-02 15:04:05-0700"),
		"homeURL":    homeURL,
	}

	if err := a.adminTpl.Execute(w, data); err != nil {
		logs.Error("app.getAdminPage:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/issue9/logs"
	"github.com/issue9/utils"
	"github.com/issue9/web"

	"github.com/caixw/gitype/helper"
)

// 将一个 log.Logger 封装成 io.Writer 接口
type logWriter log.Logger

type webhook struct {
	URL       string        `yaml:"url"`              // 接收地址，不能带域名
	Frequency time.Duration `yaml:"frequency"`        // 最小更新频率
	Method    string        `yaml:"method,omitempty"` // 请求方式，默认为 POST
	RepoURL   string        `yaml:"repoURL"`          // 远程仓库的地址
}

func (w *webhook) Sanitize() error {
	if len(w.Method) == 0 {
		w.Method = http.MethodPost
	}

	switch {
	case len(w.URL) == 0 || w.URL[0] != '/':
		return &helper.FieldError{Field: "webhook.url", Message: "不能为空且只能以 / 开头"}
	case w.Frequency < 0:
		return &helper.FieldError{Field: "webhook.frequency", Message: "不能小于 0"}
	case len(w.RepoURL) == 0:
		return &helper.FieldError{Field: "webhook.repoURL", Message: "不能为空"}
	}

	return nil
}

func (w *logWriter) Write(bs []byte) (int, error) {
	(*log.Logger)(w).Print(string(bs))
	return len(bs), nil
}

// webhooks 的回调接口
func (a *app) postWebhooks(w http.ResponseWriter, r *http.Request) {
	logs.Trace("接收到 webhook 请求")

	ctx := web.NewContext(w, r)
	if time.Now().Sub(a.client.Created()) < a.webhook.Frequency {
		logs.Error("更新过于频繁，被中止！")
		ctx.Exit(http.StatusTooManyRequests)
	}

	var cmd *exec.Cmd
	if utils.FileExists(a.path.DataDir) {
		cmd = exec.Command("git", "pull")
		cmd.Dir = a.path.DataDir
	} else {
		cmd = exec.Command("git", "clone", a.webhook.RepoURL, a.path.DataDir)
		cmd.Dir = a.path.Root
	}

	cmd.Stderr = (*logWriter)(logs.ERROR())
	cmd.Stdout = (*logWriter)(logs.INFO())
	if err := cmd.Run(); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	if err := a.reload(); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/issue9/logs"
	"github.com/issue9/utils"
)

// 将一个 log.Logger 封装成 io.Writer 接口
type logWriter log.Logger

func (w *logWriter) Write(bs []byte) (int, error) {
	(*log.Logger)(w).Print(string(bs))
	return len(bs), nil
}

// webhooks 的回调接口
func (a *app) postWebhooks(w http.ResponseWriter, r *http.Request) {
	if time.Now().Sub(a.client.Created()) < a.webhook.Frequency {
		logs.Error("更新过于频繁，被中止！")
		helper.StatusError(w, http.StatusTooManyRequests)
		return
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
		logs.Error(err)
		helper.StatusError(w, http.StatusInternalServerError)
		return
	}

	if err := a.reload(); err != nil {
		logs.Error(err)
		helper.StatusError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

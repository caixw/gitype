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
)

type logW struct {
	l *log.Logger
}

func (w *logW) Write(bs []byte) (int, error) {
	w.l.Print(string(bs))
	return len(bs), nil
}

// 通过 webhooks 来更新内容
func (a *app) postWebhooks(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Unix()

	if now-a.conf.WebhooksUpdateFreq < a.updated { // 时间太短，不接受更新
		logs.Debug("更新过于频繁，被中止！")
		return
	}

	var cmd *exec.Cmd
	if utils.FileExists(a.path.DataDir) {
		cmd = exec.Command("git", "pull")
		cmd.Dir = a.path.DataDir
	} else {
		cmd = exec.Command("git", "clone", a.conf.RepoURL, a.path.DataDir)
		cmd.Dir = a.path.Root
	}

	cmd.Stderr = &logW{l: logs.ERROR()}
	cmd.Stdout = &logW{l: logs.INFO()}
	if err := cmd.Run(); err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.reload(); err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/caixw/gitype/helper"
)

// 两个默认端口的定义
const (
	httpPort  = ":80"
	httpsPort = ":443"
)

// 对 Config.HTTPState 可选值的定义
const (
	httpStateDefault  = "default"
	httpStateDisable  = "disable"
	httpStateRedirect = "redirect"
)

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

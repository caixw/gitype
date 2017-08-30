// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/caixw/typing/data"
	"github.com/issue9/utils"
)

const httpPort = ":80"

type config struct {
	HTTPS     bool              `json:"https"`
	HTTPState string            `json:"httpState"` // 对 80 端口的处理方式，可以 disable, redirect, default
	CertFile  string            `json:"certFile"`
	KeyFile   string            `json:"keyFile"`
	Port      string            `json:"port"`
	Pprof     bool              `json:"pprof"`
	Headers   map[string]string `json:"headers"`

	Webhook       *webhook `json:"webhook"`
	AdminURL      string   `json:"adminURL"`      // 后台管理地址
	AdminPassword string   `json:"adminPassword"` // 后台管理登录地址
}

type webhook struct {
	URL       string `json:"url"`              // webhooks 接收地址
	Frequency int64  `json:"frequency"`        // webhooks 的最小更新频率，秒数
	Method    string `json:"method,omitempty"` // webhooks 的请求方式，默认为 POST
	RepoURL   string `json:"repoURL"`          // 远程仓库的地址
}

func (w *webhook) sanitize() *data.FieldError {
	if len(w.Method) == 0 {
		w.Method = http.MethodPost
	}

	switch {
	case len(w.URL) == 0 || w.URL[0] != '/':
		return &data.FieldError{Field: "webhook.URL", Message: "不能为空且只能以 / 开头"}
	case w.Frequency < 0:
		return &data.FieldError{Field: "webhook.frequency", Message: "不能小于 0"}
	case len(w.RepoURL) == 0:
		return &data.FieldError{Field: "webhook.repoURL", Message: "不能为空"}
	}

	return nil
}

func (conf *config) sanitize() *data.FieldError {
	switch {
	case conf.HTTPS && conf.HTTPState != "disable" && conf.HTTPState != "default" && conf.HTTPState != "redirect":
		return &data.FieldError{Field: "httpState", Message: "无效的取值"}
	case conf.HTTPS && conf.HTTPState != "disable" && conf.Port == httpPort:
		return &data.FieldError{Field: "port", Message: "80 端口已经被被监听"}
	case conf.HTTPS && !utils.FileExists(conf.CertFile):
		return &data.FieldError{Field: "certFile", Message: "不能为空"}
	case conf.HTTPS && !utils.FileExists(conf.KeyFile):
		return &data.FieldError{Field: "keyFile", Message: "不能为空"}
	case len(conf.AdminURL) == 0 || conf.AdminURL[0] != '/':
		return &data.FieldError{Field: "adminURL", Message: "不能为空只能以 / 开头"}
	case len(conf.AdminPassword) == 0:
		return &data.FieldError{Field: "adminPassword", Message: "不能为空"}
	}

	return conf.Webhook.sanitize()
}

func loadConfig(path string) (*config, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &config{}
	if err = json.Unmarshal(bs, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

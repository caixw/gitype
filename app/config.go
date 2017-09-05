// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

const httpPort = ":80"

const (
	httpStateDefault  = "default"
	httpStateDisable  = "disable"
	httpStateRedirect = "redirect"
)

type config struct {
	HTTPS     bool              `yaml:"https"`
	HTTPState string            `yaml:"httpState"` // 对 80 端口的处理方式，可以 disable, redirect, default
	CertFile  string            `yaml:"certFile"`
	KeyFile   string            `yaml:"keyFile"`
	Port      string            `yaml:"port"`
	Pprof     bool              `yaml:"pprof"`
	Headers   map[string]string `yaml:"headers"`
	Webhook   *webhook          `yaml:"webhook"`
}

type webhook struct {
	URL       string        `yaml:"url"`              // webhooks 接收地址
	Frequency time.Duration `yaml:"frequency"`        // webhooks 的最小更新频率
	Method    string        `yaml:"method,omitempty"` // webhooks 的请求方式，默认为 POST
	RepoURL   string        `yaml:"repoURL"`          // 远程仓库的地址
}

func loadConfig(path *vars.Path) (*config, error) {
	conf := &config{}
	if err := helper.LoadYAMLFile(path.AppConfigFile, conf); err != nil {
		return nil, err
	}

	if err := conf.sanitize(); err != nil {
		err.File = path.AppConfigFile
		return nil, err
	}

	return conf, nil
}

func (w *webhook) sanitize() *data.FieldError {
	if len(w.Method) == 0 {
		w.Method = http.MethodPost
	}

	switch {
	case len(w.URL) == 0 || w.URL[0] != '/':
		return &data.FieldError{Field: "webhook.url", Message: "不能为空且只能以 / 开头"}
	case w.Frequency < 0:
		return &data.FieldError{Field: "webhook.frequency", Message: "不能小于 0"}
	case len(w.RepoURL) == 0:
		return &data.FieldError{Field: "webhook.repoURL", Message: "不能为空"}
	}

	return nil
}

func (conf *config) sanitize() *data.FieldError {
	switch {
	case conf.HTTPS &&
		conf.HTTPState != httpStateDefault &&
		conf.HTTPState != httpStateDisable &&
		conf.HTTPState != httpStateRedirect:
		return &data.FieldError{Field: "httpState", Message: "无效的取值"}
	case conf.HTTPS &&
		conf.HTTPState != httpStateDisable &&
		conf.Port == httpPort:
		return &data.FieldError{Field: "port", Message: "80 端口已经被被监听"}
	case conf.HTTPS && !utils.FileExists(conf.CertFile):
		return &data.FieldError{Field: "certFile", Message: "不能为空"}
	case conf.HTTPS && !utils.FileExists(conf.KeyFile):
		return &data.FieldError{Field: "keyFile", Message: "不能为空"}
	}

	return conf.Webhook.sanitize()
}

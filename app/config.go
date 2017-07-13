// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"

	"github.com/caixw/typing/data"
	"github.com/issue9/utils"
)

type config struct {
	HTTPS     bool              `json:"https"`
	HTTPState string            `json:"httpState"` // 对 80 端口的处理方式，可以 disable, redirect, default
	CertFile  string            `json:"certFile"`
	KeyFile   string            `json:"keyFile"`
	Port      string            `json:"port"`
	Pprof     bool              `json:"pprof"`
	Headers   map[string]string `json:"headers"`

	WebhooksURL        string `json:"webhooksURL"`        // webhooks接收地址
	WebhooksUpdateFreq int64  `json:"webhooksUpdateFreq"` // webhooks的最小更新频率，秒数
	RepoURL            string `json:"repoURL"`            // 远程仓库的地址
	AdminURL           string `json:"adminURL"`           // 后台管理地址
	AdminPassword      string `json:"adminPassword"`      // 后台管理登录地址
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

	switch {
	case conf.HTTPS && conf.HTTPState != "disable" && conf.HTTPState != "default" && conf.HTTPState != "redirect":
		return nil, &data.FieldError{File: "config.json", Field: "httpState", Message: "无效的取值"}
	case conf.HTTPS && !utils.FileExists(conf.CertFile):
		return nil, &data.FieldError{File: "config.json", Field: "certFile", Message: "不能为空"}
	case conf.HTTPS && !utils.FileExists(conf.KeyFile):
		return nil, &data.FieldError{File: "config.json", Field: "keyFile", Message: "不能为空"}
	case len(conf.WebhooksURL) == 0 || conf.WebhooksURL[0] != '/':
		return nil, &data.FieldError{File: "config.json", Field: "webhooksURL", Message: "不能为空且只能以 / 开头"}
	case conf.WebhooksUpdateFreq < 0:
		return nil, &data.FieldError{File: "config.json", Field: "webhooksUpdateFreq", Message: "不能小于 0"}
	case len(conf.RepoURL) == 0:
		return nil, &data.FieldError{File: "config.json", Field: "repoURL", Message: "不能为空"}
	case len(conf.AdminURL) == 0 || conf.AdminURL[0] != '/':
		return nil, &data.FieldError{File: "config.json", Field: "adminURL", Message: "不能为空只能以 / 开头"}
	case len(conf.AdminPassword) == 0:
		return nil, &data.FieldError{File: "config.json", Field: "adminPassword", Message: "不能为空"}
	}

	return conf, nil
}

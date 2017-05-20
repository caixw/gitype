// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"

	"github.com/issue9/utils"
)

// 配置文件
type config struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	Port     string `json:"port"`
	Pprof    string `json:"pprof"`
	Headers  string `json:"headers"`

	WebhooksURL        string `json:"webhooksURL"`        // webhooks接收地址
	WebhooksUpdateFreq int64  `json:"webhooksUpdateFreq"` // webhooks的最小更新频率，秒数
	RepoURL            string `json:"repoURL"`            // 远程仓库的地址
	AdminURL           string `json:"adminURL"`           // 后台管理地址
	AdminPassword      string `json:"adminPassword"`      // 后台管理登录地址
}

func loadConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &config{}
	if err = json.Unmarshal(data, conf); err != nil {
		return nil, err
	}

	switch {
	case !utils.FileExists(conf.CertFile):
		return nil, &data.FieldError{File: "config.json", Field: "certFile", "不能为空"}
	case !utils.FileExists(conf.KeyFile):
		return nil, &data.FieldError{File: "config.json", Field: "keyFile", "不能为空"}
	case len(conf.Pprof) > 0 && conf.Pprof[0] != '/':
		return nil, &data.FieldError{File: "config.json", Field: "pprof", "必须以 / 开头"}
	case len(conf.WebhooksURL) == 0 || conf.WebhooksURL == "/":
		return nil, &data.FieldError{File: "config.json", Field: "webhooksURL", "不能为空且必须以 / 开头"}
	case conf.WebhooksUpdateFreq < 0:
		return nil, &data.FieldError{File: "config.json", Field: "webhooksUpdateFreq", "不能小于 0"}
	case len(conf.RepoURL) == 0 || conf.RepoURL == "/":
		return nil, &data.FieldError{File: "config.json", Field: "repoURL", "不能为空且必须以 / 开头"}
	case len(conf.AdminURL) == 0 || conf.AdminURL == "/":
		return nil, &data.FieldError{File: "config.json", Field: "adminURL", "不能为空且必须以 / 开头"}
	case len(conf.AdminPassword) == 0:
		return nil, &data.FieldError{File: "config.json", Field: "adminPassword", "不能为空"}
	}

	return conf, nil
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/issue9/handlers"
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
		return nil, errors.New("配置文件必须指定 certFile")
	case !utils.FileExists(conf.KeyFile):
		return nil, errors.New("配置文件必须指定 keyFile")
	case len(conf.Pprof) > 0 && conf.Pprof[0] != '/':
		return nil, errors.New("配置文件 pprof 必须以 / 开头")
	case len(conf.WebhooksURL) == 0 || conf.WebhooksURL == "/":
		return nil, errors.New("配置文件必须指定 webhooksURL 的值且不能为 /")
	case conf.WebhooksUpdateFreq < 0:
		return nil, errors.New("webhooksUpdateFreq 不能小于 0")
	case len(conf.RepoURL) == 0 || conf.RepoURL == "/":
		return nil, errors.New("配置文件必须指定 repoURL 的值且不能为 /")
	case len(conf.AdminURL) == 0 || conf.AdminURL == "/":
		return nil, errors.New("配置文件必须指定 adminURL 的值且不能为 /")
	case len(conf.AdminPassword) == 0:
		return nil, errors.New("配置文件必须指定 adminPassword")
	}

	if conf.isDebug() {
		conf.Core.ErrHandler = handlers.PrintDebug
	}

	return conf, nil
}

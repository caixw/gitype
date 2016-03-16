// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/issue9/handlers"
	"github.com/issue9/web"
)

// 配置文件
type config struct {
	Core               *web.Config `json:"core"`
	WebhooksURL        string      `json:"webhooksURL"`        // webhooks接收地址
	WebhooksUpdateFreq int64       `json:"webhooksUpdateFreq"` // webhooks的最小更新频率
	RepoURL            string      `json:"repoURL"`            // 远程仓库的地址
	AdminURL           string      `json:"adminURL"`           // 后台管理地址
	AdminPassword      string      `json:"adminPassword"`      // 后台管理登录地址
}

// 是否处于调试状态
func (conf *config) isDebug() bool {
	return len(conf.Core.Pprof) == 0
}

func loadConfig(path string) (*config, error) {
	// 加载程序配置
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &config{}
	if err = json.Unmarshal(data, conf); err != nil {
		return nil, err
	}

	switch {
	case len(conf.WebhooksURL) == 0 || conf.WebhooksURL == "/":
		return nil, errors.New("配置文件必须指定webhooksURL的值且不能为/")
	case conf.WebhooksUpdateFreq < 0:
		return nil, errors.New("webhooksUpdateFreq不能小于0")
	case len(conf.RepoURL) == 0 || conf.RepoURL == "/":
		return nil, errors.New("配置文件必须指定repoURL的值且不能为/")
	case len(conf.AdminURL) == 0 || conf.AdminURL == "/":
		return nil, errors.New("配置文件必须指定adminURL的值且不能为/")
	case len(conf.AdminPassword) == 0:
		return nil, errors.New("配置文件必须指定adminPassword")
	}

	if conf.isDebug() {
		conf.Core.ErrHandler = handlers.PrintDebug
	}

	return conf, nil
}

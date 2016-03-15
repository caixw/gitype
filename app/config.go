// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"

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

	if len(conf.WebhooksURL) == 0 || conf.WebhooksURL == "/" {
		return nil, errors.New("配置文件必须指定webhooksURL的值且不能为/")
	}
	if conf.WebhooksUpdateFreq < 0 {
		return nil, errors.New("webhooksUpdateFreq不能小于0")
	}
	if len(conf.RepoURL) == 0 || conf.RepoURL == "/" {
		return nil, errors.New("配置文件必须指定repoURL的值且不能为/")
	}
	if len(conf.AdminURL) == 0 || conf.AdminURL == "/" {
		return nil, errors.New("配置文件必须指定adminURL的值且不能为/")
	}
	if len(conf.AdminPassword) == 0 {
		return nil, errors.New("配置文件必须指定adminPassword")
	}

	return conf, nil
}

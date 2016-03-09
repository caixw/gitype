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

type config struct {
	Core          *web.Config `json:"core"`
	AdminURL      string      `json:"adminURL"`
	AdminPassword string      `json:"adminPassword"`
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

	if len(conf.AdminURL) == 0 || conf.AdminURL == "/" {
		return nil, errors.New("配置文件必须指定adminURL的值且不能为/")
	}
	if len(conf.AdminPassword) == 0 {
		return nil, errors.New("配置文件必须指定adminPassword")
	}

	return conf, nil
}

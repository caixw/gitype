// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/issue9/web"
)

// Config 表示程序级别的配置，修改这些配置需要重启程序才能启作用，
// 比如数据库初始化信息，路由项设置等。
type Config struct {
	Core *web.Config `json:"core"`

	DBDSN    string `json:"dbDSN"`
	DBPrefix string `json:"dbPrefix"`
	DBDriver string `json:"dbDriver"`

	FrontAPIPrefix string `json:"frontApiPrefix"` // 前端api地址前缀
	AdminAPIPrefix string `json:"adminApiPrefix"` // 后台api地址前经

	ThemeDir string `json:"themeDir"` // 主题文件所在的目录
}

// LoadConfig 用于加载path的内容，并尝试将其转换成Config实例。
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// 检测AdminAPIPrefix是否符合要求
	if len(cfg.AdminAPIPrefix) == 0 {
		return nil, errors.New("必须指定adminApiPrefix值")
	}
	if cfg.AdminAPIPrefix[len(cfg.AdminAPIPrefix)-1] == '/' {
		return nil, errors.New("adminApiPrefix不能以/符号结尾")
	}

	// 检测FrontAPIPrefix是否符合要求
	if len(cfg.FrontAPIPrefix) == 0 {
		return nil, errors.New("必须指定frontApiPrefix值")
	}
	if cfg.FrontAPIPrefix[len(cfg.FrontAPIPrefix)-1] == '/' {
		return nil, errors.New("frontApiPrefix不能以/符号结尾")
	}

	if len(cfg.ThemeDir) == 0 {
		return nil, errors.New("themeDir未指定")
	}
	if cfg.ThemeDir[len(cfg.ThemeDir)-1] != '/' {
		return nil, errors.New("themeDir只能以/结尾")
	}

	if len(cfg.DBDSN) == 0 {
		return nil, errors.New("app.json中未指定dbDSN")
	}

	if len(cfg.DBDriver) == 0 {
		return nil, errors.New("app.json中未指定dbDriver")
	}

	return cfg, nil
}

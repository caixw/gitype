// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"net/http"
	"os"
	"time"

	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/path"
	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// 输出的默认配置内容
var defaultConfig = &Config{
	HTTPS:        true,
	HTTPState:    HTTPStateRedirect,
	CertFile:     "cert",
	KeyFile:      "key",
	Port:         HTTPSPort,
	CookieMaxAge: cookieMaxAge,
	Headers: map[string]string{
		"Server": vars.Name + vars.Version(),
	},
	Webhook: &Webhook{
		URL:       "/webhooks",
		Frequency: time.Minute,
		Method:    http.MethodPost,
		RepoURL:   "https://github.com/caixw/blogs",
	},
}

// Init 初始化配置文件到 path 中
func Init(path *path.Path) error {
	if !utils.FileExists(path.ConfDir) {
		if err := os.Mkdir(path.ConfDir, os.ModePerm); err != nil {
			return err
		}
	}

	return helper.DumpYAMLFile(path.AppConfigFile, defaultConfig)
}

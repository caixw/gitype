// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package install

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/caixw/typing/core"
	"github.com/issue9/web"
)

// 输出默认的日志配置文件到指定的位置。
func OutputLogsConfigFile(path string) error {
	return ioutil.WriteFile(path, logFile, os.ModePerm)
}

// 输出配置文件到指定的位置
func OutputConfigFile(path string) error {
	cfg := &core.Config{
		Core: &web.Config{
			HTTPS:      false,
			CertFile:   "",
			KeyFile:    "",
			Port:       "8080",
			ServerName: "typing",
			Static: map[string]string{
				"/admin": "./static/admin/",
			},
		},

		DBDSN:    "./output/main.db",
		DBPrefix: "typing_",
		DBDriver: "sqlite3",

		FrontAPIPrefix: "/api",
		AdminAPIPrefix: "/admin/api",

		ThemeDir: "./static/front/",
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/caixw/typing/app/static"
	"github.com/issue9/web"
)

// 用于输出配置文件到指定的位置。
// 目前包含了日志配置文件和程序本身的配置文件。
func Install() error {
	if err := ioutil.WriteFile(logConfigPath, static.LogConfig, os.ModePerm); err != nil {
		return err
	}

	cfg := &Config{
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
		ThemeURLPrefix: "/themes",
		ThemeDir:       "./static/front/themes/",
		TempDir:        "./output/temp/",

		UploadDir:       "./output/uploads/",
		UploadDirFormat: "2006/01/",
		UploadExts:      ".txt;.png;.jpg;.jpeg",
		UploadSize:      1024 * 1024 * 5,
		UploadURLPrefix: "/uploads",
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, os.ModePerm)
}

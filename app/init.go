// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// 输出的默认配置内容
var defaultConfig = &config{
	HTTPS:     true,
	HTTPState: httpStateRedirect,
	CertFile:  "cert",
	KeyFile:   "key",
	Port:      ":443",
	Pprof:     false,
	Headers: map[string]string{
		"Server": vars.AppName + vars.Version(),
	},
	Webhook: &webhook{
		URL:       "/webhooks",
		Frequency: time.Minute,
		Method:    http.MethodPost,
		RepoURL:   "https://github.com/caixw/blogs",
	},
}

// Init 执行初始化命令
func Init(path *vars.Path) error {
	if !utils.FileExists(path.Root) {
		if err := os.Mkdir(path.Root, os.ModePerm); err != nil {
			return err
		}
	}

	if err := initConfDir(path); err != nil {
		return err
	}

	if err := data.Init(path); err != nil {
		return err
	}

	_, err := fmt.Fprintf(vars.CMDOutput, "操作成功，你现在可以在 %s 中修改具体的参数配置！\n", path.Root)
	return err
}

// 初始化 conf 目录下的数据
func initConfDir(path *vars.Path) error {
	if !utils.FileExists(path.ConfDir) {
		if err := os.Mkdir(path.ConfDir, os.ModePerm); err != nil {
			return err
		}
	}

	// logs.xml
	file, err := os.Create(path.LogsConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(defaultLogsXML); err != nil {
		return err
	}

	// app.yaml
	return helper.DumpYAMLFile(path.AppConfigFile, defaultConfig)
}

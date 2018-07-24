// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"os"
	"time"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	"github.com/issue9/utils"
	"github.com/issue9/web"
)

// 从 /testdata/conf.logs.xml 而来
var defaultLogsXML = `<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info prefix="[INFO]" flag="">
        <console output="stderr" foreground="green" background="black" />
        <rotate prefix="info-" dir="./testdata/logs/" size="5M" />
    </info>
 
    <debug prefix="[DEBUG]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="debug-" dir="./testdata/logs/" size="5M" />
    </debug>
 
    <trace prefix="[TRACE]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="trace-" dir="./testdata/logs/" size="5M" />
    </trace>
 
    <warn prefix="[WARN]">
        <console output="stderr" foreground="yellow" background="blue" />
        <rotate prefix="warn-" dir="./testdata/logs/" size="5M" />
    </warn>
 
    <error prefix="[ERROR]" flag="log.llongfile">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="error-" dir="./testdata/logs/" size="5M" />
    </error>
 
    <critical prefix="[CRITICAL]" flag="log.llongfile">
        <console output="stderr" foreground="red" background="blue" />
        <rotate prefix="critical-" dir="./testdata/logs/" size="5M" />
    </critical>
</logs>
`

// 输出的默认配置内容
var defaultConfig = &webhook{
	URL:       "/webhooks",
	Frequency: time.Minute,
	Method:    http.MethodPost,
	RepoURL:   "https://github.com/caixw/blogs",
}

// Init 初始化整个工作目录
func Init(path *path.Path) error {
	if !utils.FileExists(path.Root) {
		if err := os.Mkdir(path.Root, os.ModePerm); err != nil {
			return err
		}
	}

	if err := initConfDir(path); err != nil {
		return err
	}

	return data.Init(path)
}

// 初始化 conf 目录下的数据
func initConfDir(path *path.Path) error {
	if !utils.FileExists(path.ConfDir) {
		if err := os.Mkdir(path.ConfDir, os.ModePerm); err != nil {
			return err
		}
	}

	// logs.xml
	if err := helper.DumpTextFile(web.File("logs.xml"), defaultLogsXML); err != nil {
		return err
	}

	// app.yaml
	return helper.DumpYAMLFile(web.File("webhook.yaml"), defaultConfig)
}

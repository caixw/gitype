// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 定义了需要用到的一些路径信息
package path

import "path/filepath"

type Path struct {
	Root string // 根目录
	Data string // 数据所在目录
	Conf string // 程序配置所在目录

	ConfApp  string
	ConfLogs string

	DataConf   string
	DataTags   string
	DataURLS   string
	DataPosts  string
	DataThemes string
	DataRaws   string
}

func New(root string) *Path {
	conf := filepath.Join(root, "conf")
	data := filepath.Join(root, "data")
	meta := filepath.Join(data, "meta")

	return &Path{
		Root: root,
		Data: data,
		Conf: conf,

		ConfApp:  filepath.Join(conf, "app.json"),
		ConfLogs: filepath.Join(conf, "logs.xml"),

		DataConf:   filepath.Join(meta, "config.yaml"),
		DataTags:   filepath.Join(meta, "tags.yaml"),
		DataURLS:   filepath.Join(meta, "urls.yaml"),
		DataPosts:  filepath.Join(data, "posts"),
		DataThemes: filepath.Join(data, "themes"),
		DataRaws:   filepath.Join(data, "raws"),
	}
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import "path/filepath"

// Path 定义了程序需要用到的各种目录结构。
// 统一定义方便修改目录结构，也不容易发生因输入错误造成的BUG。
type Path struct {
	Root string // 根目录
	Data string // data 数据所在目录
	Conf string // conf 程序配置所在目录

	ConfApp  string // conf/app.json
	ConfLogs string // conf/logs.xml

	DataConf   string // data/meta/config.yaml
	DataTags   string // data/meta/tags.yaml
	DataLinks  string // data/meta/links.yaml
	DataURLS   string // data/meta/urls.yaml
	DataPosts  string // data/posts
	DataThemes string // data/themes
	DataRaws   string // data/raws
}

func NewPath(root string) (*Path, error) {
	if !filepath.IsAbs(root) {
		var err error
		root, err = filepath.Abs(root)
		if err != nil {
			return nil, err
		}
	}

	conf := filepath.Join(root, "conf")
	data := filepath.Join(root, "data")
	meta := filepath.Join(data, "meta")

	return &Path{
		Root: root,
		Data: data,
		Conf: conf,

		ConfApp:  filepath.Join(conf, "app.json"),
		ConfLogs: filepath.Join(conf, "logs.xml"),

		DataConf:  filepath.Join(meta, "config.yaml"),
		DataTags:  filepath.Join(meta, "tags.yaml"),
		DataLinks: filepath.Join(meta, "links.yaml"),
		DataURLS:  filepath.Join(meta, "urls.yaml"),

		DataPosts:  filepath.Join(data, "posts"),
		DataThemes: filepath.Join(data, "themes"),
		DataRaws:   filepath.Join(data, "raws"),
	}, nil
}

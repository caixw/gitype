// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 自定义URL
type URLS struct {
	Root   string `yaml:"root,omitempty"`   // 根地址
	Suffix string `yaml:"suffix,omitempty"` // 地址后缀
	Posts  string `yaml:"posts,omitempty"`  // 列表页地址
	Post   string `yaml:"post,omitempty"`   // 文章详细页地址
	Tags   string `yaml:"tags,omitempty"`   // 标签列表页地址
	Tag    string `yaml:"tag,omitempty"`    // 标签详细页地址
	Themes string `yaml:"themes,omitempty"` // 主题地址
	Atom   string `yaml:"atom"`             // atom 地址
	RSS    string `yaml:"rss"`              // RSS 地址
}

func (d *Data) loadURLS(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	urls := &URLS{}
	if err = yaml.Unmarshal(data, urls); err != nil {
		return err
	}
	d.URLS = urls
	return nil
}

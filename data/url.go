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
	Root   string `yaml:"root"`   // 根地址
	Suffix string `yaml:"suffix"` // 地址后缀
	Posts  string `yaml:"posts"`  // 列表页地址
	Post   string `yaml:"post"`   // 文章详细页地址
	Tags   string `yaml:"tags"`   // 标签列表页地址
	Tag    string `yaml:"tag"`    // 标签详细页地址
	Themes string `yaml:"themes"` // 主题地址
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

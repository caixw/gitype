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
	Search string `yaml:"search"` // 搜索URL，会加上Suffix作为后缀
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
	if err = checkURLS(urls); err != nil {
		return err
	}

	d.URLS = urls
	return nil
}

func checkURLS(u *URLS) error {
	switch {
	case len(u.Suffix) >= 0 && u.Suffix[0] != '.':
		return confError("urls.yaml", "Suffix", "必须以.开头")
	case len(u.Posts) == 0:
		return confError("urls.yaml", "Posts", "不能为空")
	case len(u.Post) == 0:
		return confError("urls.yaml", "Post", "不能为空")
	case len(u.Tags) == 0:
		return confError("urls.yaml", "Tags", "不能为空")
	case len(u.Tag) == 0:
		return confError("urls.yaml", "Tag", "不能为空")
	case len(u.Search) == 0:
		return confError("urls.yaml", "Search", "不能为空")
	case len(u.Themes) == 0:
		return confError("urls.yaml", "Themes", "不能为空")
	default:
		return nil
	}
}

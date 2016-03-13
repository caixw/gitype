// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 描述链接内容
type Link struct {
	Icon  string `yaml:"icon,omitempty"`  // 链接对应的图标名称，fontawesome图标名称，不用带fa-前缀。
	Title string `yaml:"title,omitempty"` // 链接的title属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text'`            // 链接的广西
}

func (d *Data) loadLinks(p string) error {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	links := make([]*Link, 0, 20)
	if err = yaml.Unmarshal(data, &links); err != nil {
		return err
	}

	if err = checkLinks("links.yaml", "", links); err != nil {
		return err
	}
	d.Links = links
	return nil
}

// 检测一组链接是否符合要求
func checkLinks(file, field string, links []*Link) error {
	for index, link := range links {
		if len(link.Text) == 0 {
			return fmt.Errorf("文件[%v]的[%v[%v]].Text错误:不能为空", file, field, index)
		}

		if len(link.URL) == 0 {
			return fmt.Errorf("文件[%v]的[%v[%v]].URL错误:不能为空", file, field, index)
		}
	}

	return nil
}

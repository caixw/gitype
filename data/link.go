// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"errors"
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

	if err = checkLinks(links); err != nil {
		return err
	}
	d.Links = links
	return nil
}

// 检测单个链接是否符合要求
func checkLink(l *Link) error {
	if len(l.Text) == 0 {
		return errors.New("未指text")
	}

	if len(l.URL) == 0 {
		return errors.New("链接未指url")
	}

	return nil
}

// 检测一组链接是否符合要求
func checkLinks(links []*Link) error {
	for index, link := range links {
		if len(link.Text) == 0 {
			return fmt.Errorf("第[%v]个链接未指text", index)
		}

		if len(link.URL) == 0 {
			return fmt.Errorf("第[%v]个链接未指url", index)
		}
	}

	return nil
}

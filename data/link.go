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
	Icon  string `yaml:"icon,omitempty"`
	Title string `yaml:"title,omitempty"`
	URL   string `yaml:"url"`
	Text  string `yaml:"text'`
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

	for index, link := range links {
		if len(link.Text) == 0 {
			return fmt.Errorf("第[%v]个链接未指text", index)
		}

		if len(link.URL) == 0 {
			return fmt.Errorf("第[%v]个链接未指url", index)
		}
	}
	d.Links = links
	return nil
}

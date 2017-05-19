// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
)

func (d *Data) loadLinks(p string) error {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	links := make([]*Link, 0, 20)
	if err = yaml.Unmarshal(data, &links); err != nil {
		return &FieldError{File: "links.yaml", Message: err.Error()}
	}

	// 检测错误
	for index, link := range links {
		if err := link.check(); err != nil {
			err.File = "links.yaml"
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	d.Links = links
	return nil
}

func (link *Link) check() *FieldError {
	if len(link.Text) == 0 {
		return &FieldError{Field: "Text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &FieldError{Field: "URL", Message: "不能为空"}
	}

	return nil
}

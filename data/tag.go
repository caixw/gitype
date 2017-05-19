// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"path"
	"strconv"

	"gopkg.in/yaml.v2"
)

func (d *Data) loadTags(p string) error {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	tags := make([]*Tag, 0, 100)
	if err = yaml.Unmarshal(data, &tags); err != nil {
		return &FieldError{File: "tags.yaml", Message: err.Error()}
	}
	for index, tag := range tags {
		if len(tag.Slug) == 0 {
			return &FieldError{File: "tags.yaml", Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Slug"}
		}

		if len(tag.Title) == 0 {
			return &FieldError{File: "tags.yaml", Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Title"}
		}

		if len(tag.Content) == 0 {
			return &FieldError{File: "tags.yaml", Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Content"}
		}

		tag.Posts = make([]*Post, 0, 10)
		tag.Permalink = path.Join(d.Config.URLS.Root, d.Config.URLS.Tag, tag.Slug+d.Config.URLS.Suffix)
	}
	d.Tags = tags
	return nil
}

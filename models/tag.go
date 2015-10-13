// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

type Tag struct {
	ID          int64  `orm:"name(id);ai" json:"id"`
	Name        string `orm:"name(name);unique(u_name);len(50);" json:"name,omitempty"` // 唯一名称
	Title       string `orm:"name(title);unique(u_title);len(50)" json:"title"`         // 名称
	Description string `orm:"name(description);len(-1)" json:"description"`             // 详细描述，可以用html
}

func (t *Tag) Meta() string {
	return `name(tags)`
}

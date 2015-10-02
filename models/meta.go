// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

const (
	MetaTypeAll = iota
	MetaTypeCat
	MetaTypeTag
)

const MetaNoParent = -1

type Meta struct {
	ID          int64  `orm:"name(id);ai"`
	Name        string `orm:"name(name);len(50);unique(unq_name);nullable" json:"name,omitempty"` // 唯一名称
	Parent      int64  `orm:"name(parent)" json:"parent,omitempty"`                               // 上级类别
	Type        int    `orm:"name(type)" json:"type,omitempty"`                                   // 类型
	Order       int    `orm:"name(order)" json:"order,omitempty"`                                 // 显示顺序
	Title       string `orm:"name(title);len(50)" json:"title"`                                   // 名称
	Description string `orm:"name(description);len(-1)" json:"description"`                       // 详细描述，可以用html
}

func (m *Meta) Meta() string {
	return `name(metas)`
}

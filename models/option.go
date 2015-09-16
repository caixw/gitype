// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

// 获取评论时的返回顺序
const (
	CommentOrderDesc = iota
	CommentOrderAsc
)

// 系统设置项。
type Option struct {
	Key   string `orm:"name(key);len(20);pk"` // 该设置项的唯一名称
	Value string `orm:"name(value);len(-1)"`  // 该设置项的值
	Group string `orm:"name(group);len(20)"`  // 对该设置项的分组。
}

func (opt *Option) Meta() string {
	return `name(options)`
}

// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

const (
	PostStateAll       = iota // 表示所有状态
	PostStatePublished        // 已经发布
	PostStateDraft            // 草稿
)

// 文章内容
type Post struct {
	ID       int64  `orm:"name(id);ai"`
	Name     string `orm:"name(name);len(200);index(idx_name)"` // 唯一名称，为空表示使用id
	Title    string `orm:"name(title);len(200)"`                // 标题
	Summary  string `orm:"name(summary);len(5000)"`             // 内容摘要
	Content  string `orm:"name(content);len(-1)"`               // 实际内容
	State    int    `orm:"name(state)"`                         // 状态
	Order    int    `orm:"name(order)"`                         // 排序
	Template string `orm:"name(template);len(50)"`              // 使用的模板

	Created  int64 `orm:"name(created)"`  // 创建时间
	Modified int64 `orm:"name(modified)"` // 最后次修改时间

	AllowPing    bool `orm:"name(allowPing)"`
	AllowComment bool `orm:"name(allowComment)"`

	//License string `orm:"name(license)"`
}

func (p *Post) Meta() string {
	return `name(posts)`
}

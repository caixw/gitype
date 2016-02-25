// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import "strconv"

// 构建绝对地址。
func (d *Data) URL(path string) string {
	return d.Config.URL + path
}

// 生成文章页URL
//  /posts/p1.html
func (d *Data) PostURL(slug string) string {
	return "/posts/" + slug + d.Config.Suffix
}

// 生成列表页URL
//  /
//  /posts.html?page=2
func (d *Data) PostsURL(page int) string {
	if page <= 1 {
		return "/"
	}

	return "/posts" + d.Config.Suffix + "?page=" + strconv.Itoa(page)
}

// 生成标签详情页URL
//  /tags/tag1.html
//  /tags/tag1.html?page=2
func (d *Data) TagURL(slug string, page int) string {
	base := "/tags/" + slug + d.Config.Suffix
	if page <= 1 {
		return base
	}

	return base + "?page=" + strconv.Itoa(page)
}

// 生成标签页URL
// /tags.html
func (d *Data) TagsURL() string {
	return "/tags" + d.Config.Suffix
}

// 生成主题页面地址
//  /themes/default/xx.css
func (d *Data) ThemeURL(path string) string {
	return "/themes" + path
}

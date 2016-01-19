// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import "strconv"

// 生成文章的url，postSlug为文章的唯一标记表示，一般为Name或是id字段。
//  /posts/about.html
func PostURL(postSlug string) string {
	return "/posts/" + postSlug + options.Suffix
}

// 为一个评论生成唯一id值
func CommentFragment(id int64) string {
	return "comments-" + strconv.FormatInt(id, 10)
}

// 生成文章评论URL，postSlug为文章的唯一标记表示，一般为Name或是id字段，id为评语的id
func CommentURL(postSlug string, id int64) string {
	return PostURL(postSlug) + "#" + CommentFragment(id)
}

// 生成标签的url，tagID为文章的唯一标记表示，一般为Name或是id字段，page为文章的页码。
//  /tags/tag1.html  // 首页
//  /tags/tag1.html?page=2 // 其它页面
func TagURL(tagID string, page int) string {
	url := "/tags/" + tagID + options.Suffix
	if page > 1 {
		url += "?page=" + strconv.Itoa(page)
	}
	return url
}

// 主页URL，可以是"/index.html"什么的，一般为"/"
func HomeURL() string {
	return "/"
}

// 生成文章列表url，首页不显示页码。
//  / 首页
//  /posts.html?page=2 // 其它页面
func PostsURL(page int) string {
	if page <= 1 {
		return "/posts" + options.Suffix
	}

	return "/posts" + options.Suffix + "?page=" + strconv.Itoa(page)
}

// 生成标签列表url，所有标签在一个页面显示，不分页。
//  /tags.html
func TagsURL() string {
	return "/tags" + options.Suffix
}

// 生成一个绝对URL，前缀为后台设置项中的SiteURL值。若未指定该值，则直接返回原值。
func URL(path string) string {
	return options.SiteURL + path
}

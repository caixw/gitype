// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import "strconv"

// 生成文章的url，postID为文章的唯一标记表示，一般为Name或是id字段。
func PostURL(opt *Options, postID string) string {
	return opt.SiteURL + "posts/" + postID + opt.Suffix
}

// 生成标签的url，tagID为文章的唯一标记表示，一般为Name或是id字段，page为文章的页码。
func TagURL(opt *Options, tagID string, page int) string {
	url := opt.SiteURL + "tags/" + tagID + opt.Suffix
	if page > 1 {
		url += "?page=" + strconv.Itoa(page)
	}
	return url
}

// 生成文章列表url，首页不显示页码。
func PostsURL(opt *Options, page int) string {
	if page <= 1 {
		return opt.SiteURL
	}

	return opt.SiteURL + "posts" + opt.Suffix + "?page=" + strconv.Itoa(page)
}

// 生成标签列表url，所有标签在一个页面显示，不分页。
func TagsURL(opt *Options) string {
	return opt.SiteURL + "tags" + opt.Suffix
}

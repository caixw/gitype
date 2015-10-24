// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

// 页面的基本信息
type PageInfo struct {
	Title       string // 网页的title值
	SiteName    string // 网站名称
	SecondTitle string // 副标题
	Canonical   string // 当前页的唯一链接
	Keywords    string // meta.keywords的值
	Description string // meta.description的值
	AppVersion  string // 当前程序的版本号
	GoVersion   string // 编译的go版本号
	Author      string // 作者名称

	PostSize    int      // 文章数量
	CommentSize int      // 评论数量
	Tags        []Anchor // 标签列表
	Topics      []Anchor // 最新评论的10条内容
	Hots        []Anchor // 评论最多的10条内容
}

type Anchor struct {
	Link  string // 链接地址
	Title string // 地址的字面文字
	Ext   string // 扩展内容，比如title,alt等，根据链接来确定
}

// 文章的详细内容
type Post struct {
	ID           int64
	Name         string
	Title        string
	Content      string
	Author       string
	Tags         []Anchor
	Comments     int    // 评论数量
	Created      int64  // 创建时间
	Modified     int64  // 修改时间
	AllowComment bool   // 是否允许评论
	Permalink    string // 文章的链接
}

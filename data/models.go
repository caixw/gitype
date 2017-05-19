// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

// Author 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

// Tag 描述标签信息
type Tag struct {
	Slug      string  `yaml:"slug"`            // 唯一名称
	Title     string  `yaml:"title"`           // 名称
	Color     string  `yaml:"color,omitempty"` // 标签颜色。若未指定，则继承父容器
	Content   string  `yaml:"content"`         // 对该标签的详细描述
	Posts     []*Post `yaml:"-"`               // 关联的文章
	Permalink string  `yaml:"-"`               // 唯一链接
}

// Post 表示文章的信息
type Post struct {
	Slug           string  `yaml:"-"`        // 唯一名称
	Title          string  `yaml:"title"`    // 标题
	Created        int64   `yaml:"-"`        // 创建时间
	Modified       int64   `yaml:"-"`        // 修改时间
	Tags           []*Tag  `yaml:"-"`        // 关联的标签
	Author         *Author `yaml:"author"`   // 作者
	Template       string  `yaml:"template"` // 使用的模板。未指定，则使用系统默认的
	Top            bool    `yaml:"top"`      // 是否置顶，多个置顶，则按时间排序
	Summary        string  `yaml:"summary"`  // 摘要
	Content        string  `yaml:"-"`        // 内容
	CreatedFormat  string  `yaml:"created"`  // 创建时间的字符串表示形式
	ModifiedFormat string  `yaml:"modified"` // 修改时间的字符串表示形式
	TagsString     string  `yaml:"tags"`     // 关联标签的列表
	Path           string  `yaml:"path"`     // 正文的文件名，相对于meta.yaml所在的目录
	Permalink      string  `yaml:"-"`        // 文章的唯一链接
}

// Link 描述链接的内容
type Link struct {
	Icon  string `yaml:"icon,omitempty"`  // 链接对应的图标名称，fontawesome 图标名称，不用带 fa- 前缀。
	Title string `yaml:"title,omitempty"` // 链接的 title 属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
}

// URLS 自定义 URL
type URLS struct {
	Root   string `yaml:"root"`   // 根地址
	Suffix string `yaml:"suffix"` // 地址后缀
	Posts  string `yaml:"posts"`  // 列表页地址
	Post   string `yaml:"post"`   // 文章详细页地址
	Tags   string `yaml:"tags"`   // 标签列表页地址
	Tag    string `yaml:"tag"`    // 标签详细页地址
	Search string `yaml:"search"` // 搜索URL，会加上Suffix作为后缀
	Themes string `yaml:"themes"` // 主题地址
}

// Config 一些基本配置项。
type Config struct {
	Title           string `yaml:"title"`                 // 网站标题
	Subtitle        string `yaml:"subtitle,omitempty"`    // 网站副标题
	URL             string `yaml:"url"`                   // 网站的地址
	Keywords        string `yaml:"keywords,omitempty"`    // 默认情况下的keyword内容
	Description     string `yaml:"description,omitempty"` // 默认情况下的descrription内容
	Beian           string `yaml:"beian,omitempty"`       // 备案号
	Uptime          int64  `yaml:"-"`                     // 上线时间，unix时间戳，由UptimeFormat转换而来
	UptimeFormat    string `yaml:"uptime"`                // 上线时间，字符串表示
	PageSize        int    `yaml:"pageSize"`              // 每页显示的数量
	LongDateFormat  string `yaml:"longDateFormat"`        // 长时间的显示格式
	ShortDateFormat string `yaml:"shortDateFormat"`       // 短时间的显示格式
	Theme           string `yaml:"theme"`                 // 默认主题

	Menus  []*Link `yaml:"menus,omitempty"` // 菜单内容
	Author *Author `yaml:"author"`          // 默认的作者信息

	// feeds
	RSS     *RSS     `yaml:"rss,omitempty"`
	Atom    *RSS     `yaml:"atom,omitempty"`
	Sitemap *Sitemap `yaml:"sitemap,omitempty"`
}

// RSS 表示 rss 或是 atom 等 feed 的信息
type RSS struct {
	Title string `yaml:"title"` // 标题
	Size  int    `yaml:"size"`  // 显示数量
	URL   string `yaml:"url"`   // 地址
}

// Sitemap 表示 sitemap 的相关配置项
type Sitemap struct {
	URL            string  `yaml:"url"`
	EnableTag      bool    `yaml:"enableTag,omitempty"`
	TagPriority    float64 `yaml:"tagPriority"`
	PostPriority   float64 `yaml:"postPriority"`
	TagChangefreq  string  `yaml:"tagChangefreq"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package data 负责加载 data 目录下的数据。
// 会调用 github.com/issue9/logs 包的内容，调用之前需要初始化该包。
package data

import (
	"fmt"

	"github.com/caixw/typing/vars"
)

// Data 结构体包含了数据目录下所有需要加载的数据内容。
type Data struct {
	path   *vars.Path
	Config *Config  // 配置内容
	Tags   []*Tag   // map 对顺序是未定的，所以使用 slice
	Links  []*Link  // 友情链接
	Posts  []*Post  // 所有的文章列表
	Themes []*Theme // 主题，使用 slice，方便排序
}

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
	Slug           string  `yaml:"-"`                  // 唯一名称
	Title          string  `yaml:"title"`              // 标题
	Created        int64   `yaml:"-"`                  // 创建时间
	Modified       int64   `yaml:"-"`                  // 修改时间
	Tags           []*Tag  `yaml:"-"`                  // 关联的标签
	Keywords       string  `yaml:"keywords,omitempty"` // meta.keywords 标签的内容，如果为空，使用 tags
	Author         *Author `yaml:"author"`             // 作者
	Template       string  `yaml:"template"`           // 使用的模板。未指定，则使用系统默认的
	Top            bool    `yaml:"top"`                // 是否置顶，多个置顶，则按时间排序
	Summary        string  `yaml:"summary"`            // 摘要
	Content        string  `yaml:"-"`                  // 内容
	CreatedFormat  string  `yaml:"created"`            // 创建时间的字符串表示形式
	ModifiedFormat string  `yaml:"modified"`           // 修改时间的字符串表示形式
	TagsString     string  `yaml:"tags"`               // 关联标签的列表
	Path           string  `yaml:"path"`               // 正文的文件名，相对于 meta.yaml 所在的目录
	Permalink      string  `yaml:"-"`                  // 文章的唯一链接
}

// Theme 表示主题信息
type Theme struct {
	ID          string  `yaml:"-"`           // 主题的唯一 ID
	Name        string  `yaml:"name"`        // 主题名称
	Version     string  `yaml:"version"`     // 主题的版本号
	Description string  `yaml:"description"` // 主题的描述信息
	Author      *Author `yaml:"author"`      // 作者
	Path        string  `yaml:"-"`           // 主题所在的目录
	Actived     bool    `yaml:"-"`           // 是否当前正在使用的主题
	Dark        bool    `yaml:"-"`           // 夜间模式
}

// Link 描述链接的内容
type Link struct {
	Icon  string `yaml:"icon,omitempty"`  // 链接对应的图标名称，fontawesome 图标名称，不用带 fa- 前缀。
	Title string `yaml:"title,omitempty"` // 链接的 title 属性
	Rel   string `yaml:"rel,omitempty"`   // 链接的 rel 属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
}

// Config 一些基本配置项。
type Config struct {
	Title           string `yaml:"title"`                 // 网站标题
	Subtitle        string `yaml:"subtitle,omitempty"`    // 网站副标题
	URL             string `yaml:"url"`                   // 网站的地址，不包含最后的斜杠
	Keywords        string `yaml:"keywords,omitempty"`    // 默认情况下的 keyword 内容
	Description     string `yaml:"description,omitempty"` // 默认情况下的 descrription 内容
	Beian           string `yaml:"beian,omitempty"`       // 备案号
	Uptime          int64  `yaml:"-"`                     // 上线时间，unix 时间戳，由 UptimeFormat 转换而来
	UptimeFormat    string `yaml:"uptime"`                // 上线时间，字符串表示
	PageSize        int    `yaml:"pageSize"`              // 每页显示的数量
	LongDateFormat  string `yaml:"longDateFormat"`        // 长时间的显示格式
	ShortDateFormat string `yaml:"shortDateFormat"`       // 短时间的显示格式
	Theme           string `yaml:"theme"`                 // 默认主题

	Author *Author `yaml:"author"` // 默认的作者信息

	Menus []*Link `yaml:"menus,omitempty"` // 菜单内容

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

// FieldError 表示加载文件出错时，具体的错误信息
type FieldError struct {
	File    string // 所在文件
	Message string // 错误信息
	Field   string // 所在的字段
}

func (err *FieldError) Error() string {
	return fmt.Sprintf("在文件 %v 中的 %v 字段发生错误： %v", err.File, err.Field, err.Message)
}

// Load 函数用于加载一份新的数据。
func Load(path *vars.Path) (*Data, error) {
	d := &Data{
		path: path,
	}

	if err := d.loadMeta(); err != nil {
		return nil, err
	}

	if err := d.loadThemes(); err != nil {
		return nil, err
	}

	found := false
	for _, theme := range d.Themes {
		if theme.ID == d.Config.Theme {
			found = true
			break
		}
	}
	if !found {
		return nil, &FieldError{File: confFile, Message: "该主题并不存在", Field: "Theme"}
	}

	// 加载文章
	if err := d.loadPosts(); err != nil {
		return nil, err
	}

	return d, nil
}

func (link *Link) check() *FieldError {
	if len(link.Text) == 0 {
		return &FieldError{Field: "Text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &FieldError{Field: "URL", Message: "不能为空"}
	}

	return nil
}

func (author *Author) check() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "Name", Message: "不能为空"}
	}

	return nil
}

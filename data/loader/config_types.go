// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import "github.com/caixw/gitype/helper"

// 归档的类型
const (
	ArchiveTypeYear  = "year"
	ArchiveTypeMonth = "month"
)

// 归档的排序方式
const (
	ArchiveOrderDesc = "desc"
	ArchiveOrderAsc  = "asc"
)

// RSS RSS 和 Atom 相关的配置项
type RSS struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Type  string `yaml:"type,omitempty"`
	Size  int    `yaml:"size"` // 显示数量
}

// Opensearch opensearch 相关的配置
type Opensearch struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type,omitempty"`
	Title string `yaml:"title,omitempty"`

	ShortName   string `yaml:"shortName"`
	Description string `yaml:"description"`
	LongName    string `yaml:"longName,omitempty"`
	Image       *Icon  `yaml:"image,omitempty"`
}

// Sitemap sitemap 相关的配置
type Sitemap struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type,omitempty"`

	XslURL     string  `yaml:"xslURL,omitempty"`    // 为 sitemap 指定一个 xsl 文件
	Priority   float64 `yaml:"priority"`            // 默认的优先级
	Changefreq string  `yaml:"changefreq"`          // 默认的更新频率
	EnableTag  bool    `yaml:"enableTag,omitempty"` // 是否将标签相关的页面写入 sitemap

	// 文章可以指定一个专门的值
	PostPriority   float64 `yaml:"postPriority"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

// Archive 存档页的配置内容
type Archive struct {
	Order  string `yaml:"order"`            // 排序方式
	Type   string `yaml:"type,omitempty"`   // 存档的分类方式，可以按年或是按月
	Format string `yaml:"format,omitempty"` // 标题的格式化字符串
}

func (rss *RSS) sanitize(conf *Config, typ string) *helper.FieldError {
	if rss.Size <= 0 {
		return &helper.FieldError{Message: "必须大于 0", Field: typ + ".Size"}
	}
	if len(rss.URL) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: typ + ".URL"}
	}

	switch typ {
	case "rss":
		rss.Type = contentTypeRSS
	case "atom":
		rss.Type = contentTypeAtom
	default:
		panic("无效的 typ 值")
	}

	if len(rss.Title) == 0 {
		rss.Title = conf.Title
	}

	return nil
}

// 检测 opensearch 取值是否正确
func (s *Opensearch) sanitize(conf *Config) *helper.FieldError {
	switch {
	case len(s.URL) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.url"}
	case len(s.ShortName) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.shortName"}
	case len(s.Description) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "opensearch.description"}
	}

	if len(s.Type) == 0 {
		s.Type = contentTypeOpensearch
	}

	if s.Image == nil && conf.Icon != nil {
		s.Image = conf.Icon
	}

	return nil
}

// 检测 sitemap 取值是否正确
func (s *Sitemap) sanitize() *helper.FieldError {
	switch {
	case len(s.URL) == 0:
		return &helper.FieldError{Message: "不能为空", Field: "sitemap.url"}
	case s.Priority > 1 || s.Priority < 0:
		return &helper.FieldError{Message: "介于[0,1]之间的浮点数", Field: "sitemap.priority"}
	case s.PostPriority > 1 || s.PostPriority < 0:
		return &helper.FieldError{Message: "介于[0,1]之间的浮点数", Field: "sitemap.postPriority"}
	case !isChangereq(s.Changefreq):
		return &helper.FieldError{Message: "取值不正确", Field: "sitemap.changefreq"}
	case !isChangereq(s.PostChangefreq):
		return &helper.FieldError{Message: "取值不正确", Field: "sitemap.postChangefreq"}
	}

	if len(s.Type) == 0 {
		s.Type = contentTypeXML
	}

	return nil
}

func (a *Archive) sanitize() *helper.FieldError {
	if len(a.Type) == 0 {
		a.Type = ArchiveTypeYear
	} else {
		if a.Type != ArchiveTypeMonth && a.Type != ArchiveTypeYear {
			return &helper.FieldError{Message: "取值不正确", Field: "archive.type"}
		}
	}

	if len(a.Order) == 0 {
		a.Order = ArchiveOrderDesc
	} else {
		if a.Order != ArchiveOrderAsc && a.Order != ArchiveOrderDesc {
			return &helper.FieldError{Message: "取值不正确", Field: "archive.order"}
		}
	}

	return nil
}

var changereqs = []string{
	"never",
	"yearly",
	"monthly",
	"weekly",
	"daily",
	"hourly",
	"always",
}

func isChangereq(val string) bool {
	for _, v := range changereqs {
		if v == val {
			return true
		}
	}
	return false
}

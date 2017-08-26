// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/vars"
	"github.com/issue9/is"
)

const confFilename = "config.yaml"

// 归档的类型
const (
	ArchiveTypeYear  = "year"
	ArchiveTypeMonth = "month"
)

// 文章是否过时的比较方式
const (
	OutdatedTypeCreated  = "created"  // 以创建时间作为对比
	OutdatedTypeModified = "modified" // 以修改时间作为对比
)

// Config 一些基本配置项。
type Config struct {
	Title           string    `yaml:"title"`                 // 网站标题
	Language        string    `yaml:"language"`              // 语言标记，比如 zh-cmn-Hans
	Subtitle        string    `yaml:"subtitle,omitempty"`    // 网站副标题
	URL             string    `yaml:"url"`                   // 网站的地址，不包含最后的斜杠，仅在生成地址时使用
	Keywords        string    `yaml:"keywords,omitempty"`    // 默认情况下的 keyword 内容
	Description     string    `yaml:"description,omitempty"` // 默认情况下的 descrription 内容
	Beian           string    `yaml:"beian,omitempty"`       // 备案号
	Uptime          int64     `yaml:"-"`                     // 上线时间，unix 时间戳，由 UptimeFormat 转换而来
	UptimeFormat    string    `yaml:"uptime"`                // 上线时间，字符串表示
	PageSize        int       `yaml:"pageSize"`              // 每页显示的数量
	LongDateFormat  string    `yaml:"longDateFormat"`        // 长时间的显示格式
	ShortDateFormat string    `yaml:"shortDateFormat"`       // 短时间的显示格式
	Theme           string    `yaml:"theme"`                 // 默认主题
	Type            string    `yaml:"type,omitempty"`        // 所有页面的 mime type 类型，默认使用 vars.ContntTypeHTML
	Icon            *Icon     `yaml:"icon,omitempty"`        // 程序默认的图标
	Menus           []*Link   `yaml:"menus,omitempty"`       // 导航菜单
	Archive         *Archive  `yaml:"archive,omitempty"`     // 归档页的配置内容
	Outdated        *Outdated `yaml:"outdated,omitempty"`    // 文章过时内容的设置

	// 一些默认值，可在各自的配置中覆盖此值
	Author  *Author `yaml:"author"`  // 默认作者信息
	License *Link   `yaml:"license"` // 默认版权信息

	// feeds
	RSS        *RSS        `yaml:"rss,omitempty"`
	Atom       *RSS        `yaml:"atom,omitempty"`
	Sitemap    *Sitemap    `yaml:"sitemap,omitempty"`
	Opensearch *Opensearch `yaml:"opensearch,omitempty"`
}

// Outdated 描述过时文章的提示信息
type Outdated struct {
	Type           string `yaml:"type"`     // 比较的类型，创建时间或是修改时间
	DurationFormat string `yaml:"duration"` // Duration 的字符中形式，用于解析，可以使用 time.Duration 字符串，比如 100h
	Duration       int64  `yaml:"-"`        // 超时的时间，秒数
	Content        string `yaml:"content"`  // 提示的内容，普通文字，不能为 html
}

// Archive 存档页的配置内容
type Archive struct {
	Type   string `yaml:"type"`   // 存档的分类方式，可以为 year 或是 month
	Format string `yaml:"format"` // 标题的格式化字符串
}

// Opensearch 相关定义
type Opensearch struct {
	URL   string `yaml:"url"`             // opensearch 的地址，不能包含域名
	Title string `yaml:"title,omitempty"` // 出现于 html>head>link.title 属性中

	ShortName   string `yaml:"shortName"`
	Description string `yaml:"description"`
	LongName    string `yaml:"longName,omitempty"`
	Type        string `yaml:"type,omitempty"` // mimeType 默认取 vars.ContentTypeOpensearch
	Image       *Icon  `yaml:"image,omitempty"`
}

// RSS 表示 rss 或是 atom 等 feed 的信息
type RSS struct {
	Title string `yaml:"title"`          // 标题
	Size  int    `yaml:"size"`           // 显示数量
	URL   string `yaml:"url"`            // 地址，不能包含域名
	Type  string `yaml:"type,omitempty"` // mimeType
}

// Sitemap 表示 sitemap 的相关配置项
type Sitemap struct {
	URL        string  `yaml:"url"`                 // 展示给用户的地址，不能包含域名
	XslURL     string  `yaml:"xslURL,omitempty"`    // 为 sitemap 指定一个 xsl 文件
	Priority   float64 `yaml:"priority"`            // 默认的优先级
	Changefreq string  `yaml:"changefreq"`          // 默认的更新频率
	Type       string  `yaml:"type,omitempty"`      // mimeType
	EnableTag  bool    `yaml:"enableTag,omitempty"` // 是否将标签相关的页面写入 sitemap

	// 文章可以指定一个专门的值
	PostPriority   float64 `yaml:"postPriority"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

func (conf *Config) sanitize() *FieldError {
	if conf.PageSize <= 0 {
		return &FieldError{File: confFilename, Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "LongDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "ShortDateFormat"}
	}

	t, err := parseDate(conf.UptimeFormat)
	if err != nil {
		return &FieldError{File: confFilename, Message: err.Error(), Field: "UptimeFormat"}
	}
	conf.Uptime = t

	if len(conf.Type) == 0 {
		conf.Type = vars.ContentTypeHTML
	}

	// icon
	if conf.Icon != nil {
		if err := conf.Icon.sanitize(); err != nil {
			err.File = confFilename
			err.Field = "Icon." + err.Field
			return err
		}
	}

	// Author
	if conf.Author == nil {
		return &FieldError{File: confFilename, Message: "必须指定作者", Field: "Author"}
	}
	if len(conf.Author.Name) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "Author.Name"}
	}

	if len(conf.Title) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "Title"}
	}

	if !is.URL(conf.URL) {
		return &FieldError{File: confFilename, Message: "不是一个合法的域名或 IP", Field: "URL"}
	}
	if strings.HasSuffix(conf.URL, "/") {
		conf.URL = conf.URL[:len(conf.URL)-1]
	}

	// theme
	if len(conf.Theme) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "Theme"}
	}

	// archive
	if conf.Archive == nil {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "archive"}
	}
	if conf.Archive.Type != ArchiveTypeMonth && conf.Archive.Type != ArchiveTypeYear {
		return &FieldError{File: confFilename, Message: "取值不正确", Field: "archive.type"}
	}

	// outdated
	if conf.Outdated != nil {
		if err := conf.Outdated.sanitize(); err != nil {
			return err
		}
	}

	// license
	if conf.License == nil {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "license"}
	}
	if err := conf.License.sanitize(); err != nil {
		return err
	}

	// rss
	if err := checkRSS("RSS", conf.RSS); err != nil {
		return err
	}
	if conf.RSS != nil && len(conf.RSS.Title) == 0 {
		conf.RSS.Title = conf.Title
	}
	if conf.RSS != nil && len(conf.RSS.Type) == 0 {
		conf.RSS.Type = vars.ContentTypeRSS
	}

	// atom
	if err := checkRSS("Atom", conf.Atom); err != nil {
		return err
	}
	if conf.Atom != nil && len(conf.Atom.Title) == 0 {
		conf.Atom.Title = conf.Title
	}
	if conf.Atom != nil && len(conf.Atom.Type) == 0 {
		conf.Atom.Type = vars.ContentTypeAtom
	}

	// sitemap
	if err := checkSitemap(conf.Sitemap); err != nil {
		return err
	}
	if conf.Sitemap != nil && len(conf.Sitemap.Type) == 0 {
		conf.Sitemap.Type = vars.ContentTypeXML
	}

	// opensearch
	if err := checkOpensearch(conf.Opensearch); err != nil {
		return err
	}
	if conf.Opensearch != nil && len(conf.Opensearch.Type) == 0 {
		conf.Opensearch.Type = vars.ContentTypeOpensearch
	}
	if conf.Opensearch != nil &&
		conf.Opensearch.Image == nil &&
		conf.Icon != nil {
		conf.Opensearch.Image = conf.Icon
	}

	// Menus
	for index, link := range conf.Menus {
		if err := link.sanitize(); err != nil {
			err.File = confFilename
			err.Field = "Menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

func (o *Outdated) sanitize() *FieldError {
	if o.Type != OutdatedTypeCreated && o.Type != OutdatedTypeModified {
		return &FieldError{File: confFilename, Message: "无效的值", Field: "Outdated.Type"}
	}

	if len(o.Content) == 0 {
		return &FieldError{File: confFilename, Message: "不能为空", Field: "Outdated.Content"}
	}

	dur, err := time.ParseDuration(o.DurationFormat)
	if err != nil {
		return &FieldError{File: confFilename, Message: err.Error(), Field: "Outdated.Duration"}
	}
	o.Duration = int64(dur.Seconds())

	return nil
}

// 检测 RSS 是否正常
func checkRSS(typ string, rss *RSS) *FieldError {
	if rss != nil {
		if rss.Size <= 0 {
			return &FieldError{File: confFilename, Message: "必须大于 0", Field: typ + ".Size"}
		}
		if len(rss.URL) == 0 {
			return &FieldError{File: confFilename, Message: "不能为空", Field: typ + ".URL"}
		}
	}

	return nil
}

// 检测 sitemap 取值是否正确
func checkSitemap(s *Sitemap) *FieldError {
	if s != nil {
		switch {
		case len(s.URL) == 0:
			return &FieldError{File: confFilename, Message: "不能为空", Field: "Sitemap.URL"}
		case s.Priority > 1 || s.Priority < 0:
			return &FieldError{File: confFilename, Message: "介于[0,1]之间的浮点数", Field: "Sitemap.priority"}
		case s.PostPriority > 1 || s.PostPriority < 0:
			return &FieldError{File: confFilename, Message: "介于[0,1]之间的浮点数", Field: "Sitemap.PostPriority"}
		case !isChangereq(s.Changefreq):
			return &FieldError{File: confFilename, Message: "取值不正确", Field: "Sitemap.changefreq"}
		case !isChangereq(s.PostChangefreq):
			return &FieldError{File: confFilename, Message: "取值不正确", Field: "Sitemap.PostChangefreq"}
		}
	}
	return nil
}

// 检测 opensearch 取值是否正确
func checkOpensearch(s *Opensearch) *FieldError {
	if s != nil {
		switch {
		case len(s.URL) == 0:
			return &FieldError{File: confFilename, Message: "不能为空", Field: "Opensearch.URL"}
		case len(s.ShortName) == 0:
			return &FieldError{File: confFilename, Message: "不能为空", Field: "Opensearch.ShortName"}
		case len(s.Description) == 0:
			return &FieldError{File: confFilename, Message: "不能为空", Field: "Opensearch.Description"}
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

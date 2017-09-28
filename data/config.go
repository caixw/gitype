// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
	"github.com/issue9/is"
)

const contentTypeHTML = "text/html"

// Config 配置信息，相对于 config，去掉部分临时性的变量，供 Data 对外公开用。
type Config struct {
	Title           string    // 网站标题
	Language        string    // 语言标记，比如 zh-cmn-Hans
	Subtitle        string    // 网站副标题
	URL             string    // 网站的域名，非默认端口也得包含，不包含最后的斜杠，仅在生成地址时使用
	Keywords        string    // 默认情况下的 keyword 内容
	Description     string    // 默认情况下的 descrription 内容
	Beian           string    // 备案号
	Uptime          time.Time // 上线时间
	PageSize        int       // 每页显示的数量
	LongDateFormat  string    // 长时间的显示格式
	ShortDateFormat string    // 短时间的显示格式
	Type            string    // 所有页面的 mime type 类型，默认使用
	Icon            *Icon     // 程序默认的图标
	Menus           []*Link   // 导航菜单
	Author          *Author   // 默认作者信息
	License         *Link     // 默认版权信息
}

// 配置信息，用于从文件中读取
type config struct {
	Title           string    `yaml:"title"`
	Language        string    `yaml:"language"`
	Subtitle        string    `yaml:"subtitle,omitempty"`
	URL             string    `yaml:"url"`
	Keywords        string    `yaml:"keywords,omitempty"`
	Description     string    `yaml:"description,omitempty"`
	Beian           string    `yaml:"beian,omitempty"`
	Uptime          time.Time `yaml:"-"` // 上线时间，unix 时间戳，由 UptimeFormat 转换而来
	PageSize        int       `yaml:"pageSize"`
	LongDateFormat  string    `yaml:"longDateFormat"`
	ShortDateFormat string    `yaml:"shortDateFormat"`
	Type            string    `yaml:"type,omitempty"`
	Icon            *Icon     `yaml:"icon,omitempty"`
	Menus           []*Link   `yaml:"menus,omitempty"`
	Author          *Author   `yaml:"author"`
	License         *Link     `yaml:"license"`

	// 以下内容不存在于 Config 中
	Outdated     *outdatedConfig   `yaml:"outdated,omitempty"`
	Theme        string            `yaml:"theme"`
	UptimeFormat string            `yaml:"uptime"`
	Archive      *archiveConfig    `yaml:"archive"`
	RSS          *rssConfig        `yaml:"rss,omitempty"`
	Atom         *rssConfig        `yaml:"atom,omitempty"`
	Sitemap      *sitemapConfig    `yaml:"sitemap,omitempty"`
	Opensearch   *opensearchConfig `yaml:"opensearch,omitempty"`
}

func newConfig(conf *config) *Config {
	return &Config{
		Title:           conf.Title,
		Language:        conf.Language,
		Subtitle:        conf.Subtitle,
		URL:             conf.URL,
		Keywords:        conf.Keywords,
		Description:     conf.Description,
		Beian:           conf.Beian,
		Uptime:          conf.Uptime,
		PageSize:        conf.PageSize,
		LongDateFormat:  conf.LongDateFormat,
		ShortDateFormat: conf.ShortDateFormat,
		Type:            conf.Type,
		Icon:            conf.Icon,
		Menus:           conf.Menus,
	}
}

func loadConfig(path *vars.Path) (*config, error) {
	conf := &config{}
	if err := helper.LoadYAMLFile(path.MetaConfigFile, conf); err != nil {
		return nil, err
	}

	if err := conf.sanitize(); err != nil {
		err.File = path.MetaConfigFile
		return nil, err
	}

	return conf, nil
}

func (conf *config) sanitize() *helper.FieldError {
	if conf.PageSize <= 0 {
		return &helper.FieldError{Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "longDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "shortDateFormat"}
	}

	t, err := time.Parse(vars.DateFormat, conf.UptimeFormat)
	if err != nil {
		return &helper.FieldError{Message: err.Error(), Field: "uptimeFormat"}
	}
	conf.Uptime = t

	if len(conf.Type) == 0 {
		conf.Type = contentTypeHTML
	}

	// icon
	if conf.Icon != nil {
		if err := conf.Icon.sanitize(); err != nil {
			err.Field = "icon." + err.Field
			return err
		}
	}

	// Author
	if conf.Author == nil {
		return &helper.FieldError{Message: "必须指定作者", Field: "author"}
	}
	if err := conf.Author.sanitize(); err != nil {
		err.Field = "author." + err.Field
		return err
	}

	if len(conf.Title) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "title"}
	}

	if !is.URL(conf.URL) {
		return &helper.FieldError{Message: "不是一个合法的域名或 IP", Field: "url"}
	}
	if strings.HasSuffix(conf.URL, "/") {
		conf.URL = conf.URL[:len(conf.URL)-1]
	}

	// theme
	if len(conf.Theme) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "theme"}
	}

	// archive
	if conf.Archive == nil {
		return &helper.FieldError{Message: "不能为空", Field: "archive"}
	}
	if err := conf.Archive.sanitize(); err != nil {
		return err
	}

	// outdated
	if conf.Outdated != nil {
		if err := conf.Outdated.sanitize(); err != nil {
			return err
		}
	}

	// license
	if conf.License == nil {
		return &helper.FieldError{Message: "不能为空", Field: "license"}
	}
	if err := conf.License.sanitize(); err != nil {
		err.Field = "license." + err.Field
		return err
	}

	// rss
	if conf.RSS != nil {
		if err := conf.RSS.sanitize(conf, "rss"); err != nil {
			return err
		}
	}

	// atom
	if conf.Atom != nil {
		if err := conf.Atom.sanitize(conf, "atom"); err != nil {
			return err
		}
	}

	// sitemap
	if conf.Sitemap != nil {
		if err := conf.Sitemap.sanitize(); err != nil {
			return err
		}
	}

	// opensearch，需要用到 config.Icon 变量
	if conf.Opensearch != nil {
		if err := conf.Opensearch.sanitize(conf); err != nil {
			return err
		}
	}

	// Menus
	for index, link := range conf.Menus {
		if err := link.sanitize(); err != nil {
			err.Field = "Menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

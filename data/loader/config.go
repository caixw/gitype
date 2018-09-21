// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"strconv"
	"time"

	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/path"
	l "golang.org/x/text/language"
)

// 默认的语言，在配置文件中未指定时，使用此值，
// 作为默认值，此值最好不要修改，若需要修改，
// 则最好将诸如 tagTitle 等与语言相关的常量一起修改。
const language = "zh-cmn-Hans"

// Config 配置信息，用于从文件中读取
type Config struct {
	Title           string        `yaml:"title"`
	TitleSeparator  string        `yaml:"titleSeparator"`
	Language        string        `yaml:"language"`
	Subtitle        string        `yaml:"subtitle,omitempty"`
	Beian           string        `yaml:"beian,omitempty"`
	Uptime          time.Time     `yaml:"uptime"`
	PageSize        int           `yaml:"pageSize"`
	Type            string        `yaml:"type,omitempty"`
	Icon            *Icon         `yaml:"icon,omitempty"`
	Menus           []*Link       `yaml:"menus,omitempty"`
	Author          *Author       `yaml:"author"`
	License         *Link         `yaml:"license"`
	LongDateFormat  string        `yaml:"longDateFormat"`
	ShortDateFormat string        `yaml:"shortDateFormat"`
	Outdated        time.Duration `yaml:"outdated,omitempty"`
	Theme           string        `yaml:"theme"`

	// 各个页面的一些自定义项，目前支持以下几个元素的修改：
	// 1) html>head>title
	// 2) html>head>meta.keywords
	// 3) html>head>meta.description
	Pages map[string]*Page `yaml:"pages,omitempty"`

	Archive    *Archive    `yaml:"archive"`
	RSS        *RSS        `yaml:"rss,omitempty"`
	Atom       *RSS        `yaml:"atom,omitempty"`
	Sitemap    *Sitemap    `yaml:"sitemap,omitempty"`
	Opensearch *Opensearch `yaml:"opensearch,omitempty"`

	LanguageTag l.Tag `yaml:"-"`
}

// LoadConfig 加载配置信息
func LoadConfig(path *path.Path) (*Config, error) {
	conf := &Config{}
	if err := helper.LoadYAMLFile(path.MetaConfigFile, conf); err != nil {
		return nil, err
	}

	if err := conf.sanitize(); err != nil {
		err.File = path.MetaConfigFile
		return nil, err
	}

	return conf, nil
}

func (conf *Config) sanitize() *helper.FieldError {
	if len(conf.Language) == 0 {
		conf.Language = language
	}

	tag, err := l.Parse(conf.Language)
	if err != nil {
		return &helper.FieldError{Message: err.Error(), Field: "language"}
	}
	conf.LanguageTag = tag

	if conf.PageSize <= 0 {
		return &helper.FieldError{Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "longDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &helper.FieldError{Message: "不能为空", Field: "shortDateFormat"}
	}

	if conf.Outdated < 0 {
		return &helper.FieldError{Message: "必须大于 0", Field: "outdated"}
	}

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

	// menus
	for index, link := range conf.Menus {
		if err := link.sanitize(); err != nil {
			err.Field = "Menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	conf.initPages()

	return nil
}

// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/is"
	"gopkg.in/yaml.v2"
)

// 一些基本配置项。
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

type RSS struct {
	Title string `yaml:"title"` // 标题
	Size  int    `yaml:"size"`  // 显示数量
	URL   string `yaml:"url"`   // 地址
}

type Sitemap struct {
	URL            string  `yaml:"url"`
	EnableTag      bool    `yaml:"enableTag,omitempty"`
	TagPriority    float64 `yaml:"tagPriority"`
	PostPriority   float64 `yaml:"postPriority"`
	TagChangefreq  string  `yaml:"tagChangefreq"`
	PostChangefreq string  `yaml:"postChangefreq"`
}

// 加载配置文件。
// path 配置文件的地址。
func (d *Data) loadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config := &Config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		return err
	}

	// 检测变量是否正确
	if err = initConfig(config); err != nil {
		return err
	}

	d.Config = config
	return nil
}

// initConfig 初始化config的内容，负责检测数据的合法性和格式的转换。
func initConfig(conf *Config) error {
	if conf.PageSize <= 0 {
		return &MetaError{File: "config.yaml", Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &MetaError{File: "config.yaml", Message: "不能为空", Field: "LongDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &MetaError{File: "config.yaml", Message: "不能为空", Field: "ShortDateFormat"}
	}

	t, err := time.Parse(parseDateFormat, conf.UptimeFormat)
	if err != nil {
		return &MetaError{File: "config.yaml", Message: err.Error(), Field: "UptimeFormat"}
	}
	conf.Uptime = t.Unix()

	// Author
	if conf.Author == nil {
		return &MetaError{File: "config.yaml", Message: "必须指定作者", Field: "Author"}
	}
	if len(conf.Author.Name) == 0 {
		return &MetaError{File: "config.yaml", Message: "不能为空", Field: "Author.Name"}
	}

	if len(conf.Title) == 0 {
		return &MetaError{File: "config.yaml", Message: "不能为空", Field: "Title"}
	}

	if !is.URL(conf.URL) {
		return &MetaError{File: "config.yaml", Message: "不是一个合法的域名或IP", Field: "URL"}
	}
	if strings.HasSuffix(conf.URL, "/") {
		conf.URL = conf.URL[:len(conf.URL)-1]
	}

	// theme
	if len(conf.Theme) == 0 {
		return &MetaError{File: "config.yaml", Message: "不能为空", Field: "Theme"}
	}

	if err := checkRSS("RSS", conf.RSS); err != nil {
		return err
	}

	if err := checkRSS("Atom", conf.Atom); err != nil {
		return err
	}

	if err := checkSitemap(conf.Sitemap); err != nil {
		return err
	}

	// Menus
	for index, link := range conf.Menus {
		if err := link.check(); err != nil {
			err.File = "config.yaml"
			err.Field = "Menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

// 检测RSS是否正常
func checkRSS(typ string, rss *RSS) error {
	if rss != nil {
		if len(rss.Title) == 0 {
			return &MetaError{File: "config.yaml", Message: "不能为空", Field: typ + ".Title"}
		}
		if rss.Size <= 0 {
			return &MetaError{File: "config.yaml", Message: "必须大于0", Field: typ + ".Size"}
		}
		if len(rss.URL) == 0 {
			return &MetaError{File: "config.yaml", Message: "不能为空", Field: typ + ".URL"}
		}
	}

	return nil
}

// 检测sitemap取值是否正确
func checkSitemap(s *Sitemap) error {
	if s != nil {
		switch {
		case len(s.URL) == 0:
			return &MetaError{File: "config.yaml", Message: "不能为空", Field: "Sitemap.URL"}
		case s.TagPriority > 1 || s.TagPriority < 0:
			return &MetaError{File: "config.yaml", Message: "介于[0,1]之间的浮点数", Field: "Sitemap.TagPriority"}
		case s.PostPriority > 1 || s.PostPriority < 0:
			return &MetaError{File: "config.yaml", Message: "介于[0,1]之间的浮点数", Field: "Sitemap.PostPriority"}
		case !isChangereq(s.TagChangefreq):
			return &MetaError{File: "config.yaml", Message: "取值不正确", Field: "Sitemap.TagChangefreq"}
		case !isChangereq(s.PostChangefreq):
			return &MetaError{File: "config.yaml", Message: "取值不正确", Field: "Sitemap.PostChangefreq"}
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

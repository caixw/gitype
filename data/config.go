// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/issue9/is"
	"github.com/issue9/utils"
	"gopkg.in/yaml.v2"
)

// 一些基本配置项。
type Config struct {
	Logo            string `yaml:"logo,omitempty"`            // 网站的logo
	Title           string `yaml:"title"`                     // 网站标题
	Subtitle        string `yaml:"subtitle,omitempty"`        // 网站副标题
	URL             string `yaml:"url"`                       // 网站的地址
	Keywords        string `yaml:"keywords,omitempty"`        // 默认情况下的keyword内容
	Description     string `yaml:"description,omitempty"`     // 默认情况下的descrription内容
	Beian           string `yaml:"beian,omitempty"`           // 备案号
	Uptime          int64  `yaml:"-"`                         // 上线时间，unix时间戳，由UptimeFormat转换而来
	UptimeFormat    string `yaml:"uptime"`                    // 上线时间，字符串表示
	PageSize        int    `yaml:"pagesize,omitempty"`        // 每页显示的数量
	LongDateFormat  string `yaml:"longDateFormat,omitempty"`  // 长时间的显示格式
	ShortDateFormat string `yaml:"shortDateFormat,omitempty"` // 短时间的显示格式
	Theme           string `yaml:"theme,omitempty"`           // 默认主题

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

// 生成一条config字段值错误的error实例
func confError(field, message string) error {
	return fmt.Errorf("字段[%v]错误:[%v]", field, message)
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

	// 合并默认值,TODO 以安装的方式提供默认数据，而不是采有合并的方式
	conf := &Config{
		PageSize:        20,
		LongDateFormat:  "2006-01-02 15:04:05",
		ShortDateFormat: "2006-01-02",
	}
	if err = utils.Merge(true, conf, config); err != nil {
		return err
	}

	// 检测变量是否正确
	if err = checkConfig(conf, d.path.Data); err != nil {
		return err
	}

	// 做一些修正，比如时间格式转换成int64等。
	if err = fixedConfig(conf); err != nil {
		return err
	}

	d.Config = conf
	return nil
}

// 对Config实例做一些修正，比如时间格式转换成int64等。
func fixedConfig(conf *Config) error {
	// 时间转换
	if len(conf.UptimeFormat) == 0 {
		conf.Uptime = 0
	} else {
		t, err := time.Parse(parseDateFormat, conf.UptimeFormat)
		if err != nil {
			return confError("UptimeFormat", err.Error())
		}
		conf.Uptime = t.Unix()
	}

	// 确保conf.URL不能/结尾
	if strings.HasSuffix(conf.URL, "/") {
		conf.URL = conf.URL[:len(conf.URL)-1]
	}

	return nil
}

// 检测config所有变量是否合法。不合法返回eror
// path data的路径名。
func checkConfig(conf *Config, path string) error {
	if conf.PageSize <= 0 {
		return confError("pageSize", "必须为大于零的整数")
	}

	// Author
	if conf.Author == nil {
		return confError("Author", "必须指定作者")
	}
	if len(conf.Author.Name) == 0 {
		return confError("Author.Name", "不能为空")
	}

	if len(conf.Title) == 0 {
		return confError("Title", "不能为空")
	}

	if !is.URL(conf.URL) {
		return confError("URL", "不是一个合法的域名或IP")
	}

	// theme
	if len(conf.Theme) == 0 {
		return confError("Theme", "不能为空")
	}
	themes, err := getThemesName(filepath.Join(path, "themes"))
	if err != nil {
		return err
	}
	found := false
	for _, theme := range themes {
		if theme == conf.Theme {
			found = true
			break
		}
	}
	if !found {
		return confError("Theme", "该主题并不存在")
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

	return nil
}

// 检测RSS是否正常
func checkRSS(typ string, rss *RSS) error {
	if rss != nil {
		if len(rss.Title) == 0 {
			return confError(typ+".Title", "不能为空")
		}
		if rss.Size <= 0 {
			return confError(typ+".Size", "必须大于0")
		}
		if len(rss.URL) == 0 {
			return confError(typ+".URL", "不能为空")
		}
	}

	return nil
}

// 检测sitemap取值是否正确
func checkSitemap(s *Sitemap) error {
	if s != nil {
		switch {
		case len(s.URL) == 0:
			return confError("Sitemap.URL", "不能为空")
		case s.TagPriority > 1 || s.TagPriority < 0:
			return confError("Sitemap.TagPriority", "介于[0,1]之间的浮点数")
		case s.PostPriority > 1 || s.PostPriority < 0:
			return confError("Sitemap.PostPriority", "介于[0,1]之间的浮点数")
		case !isChangereq(s.TagChangefreq):
			return confError("Sitemap.TagChangefreq", "取值不正确")
		case !isChangereq(s.PostChangefreq):
			return confError("Sitemap.PostChangefreq", "取值不正确")
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

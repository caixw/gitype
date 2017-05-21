// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/caixw/typing/vars"
	"github.com/issue9/is"
	"gopkg.in/yaml.v2"
)

const (
	tagsFile  = "tags.yaml"
	confFile  = "config.yaml"
	linksFile = "links.yaml"
)

func (d *Data) loadMeta() error {
	if err := d.loadTags(); err != nil {
		return err
	}

	if err := d.loadLinks(); err != nil {
		return err
	}

	return d.loadConfig()
}

func (d *Data) metaPath(file string) string {
	return filepath.Join(d.Root, metaDir, file)
}

func (d *Data) loadTags() error {
	data, err := ioutil.ReadFile(d.metaPath(tagsFile))
	if err != nil {
		return err
	}

	tags := make([]*Tag, 0, 100)
	if err = yaml.Unmarshal(data, &tags); err != nil {
		return &FieldError{File: tagsFile, Message: err.Error()}
	}
	for index, tag := range tags {
		if len(tag.Slug) == 0 {
			return &FieldError{File: tagsFile, Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Slug"}
		}

		if len(tag.Title) == 0 {
			return &FieldError{File: tagsFile, Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Title"}
		}

		if len(tag.Content) == 0 {
			return &FieldError{File: tagsFile, Message: "不能为空", Field: "[" + strconv.Itoa(index) + "].Content"}
		}

		tag.Posts = make([]*Post, 0, 10)
		tag.Permalink = path.Join(d.Config.URLS.Root, d.Config.URLS.Tag, tag.Slug+d.Config.URLS.Suffix)
	}
	d.Tags = tags
	return nil
}

func (d *Data) loadLinks() error {
	data, err := ioutil.ReadFile(d.metaPath(linksFile))
	if err != nil {
		return err
	}

	links := make([]*Link, 0, 20)
	if err = yaml.Unmarshal(data, &links); err != nil {
		return &FieldError{File: linksFile, Message: err.Error()}
	}

	// 检测错误
	for index, link := range links {
		if err := link.check(); err != nil {
			err.File = linksFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	d.Links = links
	return nil
}

// 加载配置文件。
// path 配置文件的地址。
func (d *Data) loadConfig() error {
	data, err := ioutil.ReadFile(d.metaPath(confFile))
	if err != nil {
		return err
	}

	config := &Config{}
	if err = yaml.Unmarshal(data, config); err != nil {
		return &FieldError{File: confFile, Message: err.Error()}
	}

	// 检测变量是否正确
	if err = initConfig(config); err != nil {
		return err
	}

	d.Config = config
	return nil
}

// initConfig 初始化 config 的内容，负责检测数据的合法性和格式的转换。
func initConfig(conf *Config) error {
	if conf.PageSize <= 0 {
		return &FieldError{File: confFile, Message: "必须为大于零的整数", Field: "pageSize"}
	}

	if len(conf.LongDateFormat) == 0 {
		return &FieldError{File: confFile, Message: "不能为空", Field: "LongDateFormat"}
	}

	if len(conf.ShortDateFormat) == 0 {
		return &FieldError{File: confFile, Message: "不能为空", Field: "ShortDateFormat"}
	}

	t, err := time.Parse(vars.DateFormat, conf.UptimeFormat)
	if err != nil {
		return &FieldError{File: confFile, Message: err.Error(), Field: "UptimeFormat"}
	}
	conf.Uptime = t.Unix()

	// Author
	if conf.Author == nil {
		return &FieldError{File: confFile, Message: "必须指定作者", Field: "Author"}
	}
	if len(conf.Author.Name) == 0 {
		return &FieldError{File: confFile, Message: "不能为空", Field: "Author.Name"}
	}

	if len(conf.Title) == 0 {
		return &FieldError{File: confFile, Message: "不能为空", Field: "Title"}
	}

	if !is.URL(conf.URL) {
		return &FieldError{File: confFile, Message: "不是一个合法的域名或 IP", Field: "URL"}
	}
	if strings.HasSuffix(conf.URL, "/") {
		conf.URL = conf.URL[:len(conf.URL)-1]
	}

	// theme
	if len(conf.Theme) == 0 {
		return &FieldError{File: confFile, Message: "不能为空", Field: "Theme"}
	}

	if err := checkRSS("RSS", conf.RSS); err != nil {
		return err
	}
	if conf.RSS != nil && len(conf.RSS.Title) == 0 {
		conf.RSS.Title = conf.Title
	}

	if err := checkRSS("Atom", conf.Atom); err != nil {
		return err
	}
	if conf.Atom != nil && len(conf.Atom.Title) == 0 {
		conf.Atom.Title = conf.Title
	}

	if err := checkSitemap(conf.Sitemap); err != nil {
		return err
	}

	if conf.URLS == nil {
		return &FieldError{File: confFile, Field: "urls", Message: "不能为空"}
	}
	if err := checkURLS(conf.URLS); err != nil {
		return err
	}

	// Menus
	for index, link := range conf.Menus {
		if err := link.check(); err != nil {
			err.File = confFile
			err.Field = "Menus[" + strconv.Itoa(index) + "]." + err.Field
			return err
		}
	}

	return nil
}

// 检测 RSS 是否正常
func checkRSS(typ string, rss *RSS) error {
	if rss != nil {
		if rss.Size <= 0 {
			return &FieldError{File: confFile, Message: "必须大于0", Field: typ + ".Size"}
		}
		if len(rss.URL) == 0 {
			return &FieldError{File: confFile, Message: "不能为空", Field: typ + ".URL"}
		}
	}

	return nil
}

func checkURLS(u *URLS) error {
	switch {
	case len(u.Suffix) >= 0 && u.Suffix[0] != '.':
		return &FieldError{File: confFile, Field: "Suffix", Message: "必须以.开头"}
	case len(u.Posts) == 0:
		return &FieldError{File: confFile, Field: "Posts", Message: "不能为空"}
	case len(u.Media) == 0:
		return &FieldError{File: confFile, Field: "Media", Message: "不能为空"}
	case len(u.Post) == 0:
		return &FieldError{File: confFile, Field: "Post", Message: "不能为空"}
	case len(u.Tags) == 0:
		return &FieldError{File: confFile, Field: "Tags", Message: "不能为空"}
	case len(u.Tag) == 0:
		return &FieldError{File: confFile, Field: "Tag", Message: "不能为空"}
	case len(u.Search) == 0:
		return &FieldError{File: confFile, Field: "Search", Message: "不能为空"}
	case len(u.Themes) == 0:
		return &FieldError{File: confFile, Field: "Themes", Message: "不能为空"}
	default:
		return nil
	}
}

// 检测 sitemap 取值是否正确
func checkSitemap(s *Sitemap) error {
	if s != nil {
		switch {
		case len(s.URL) == 0:
			return &FieldError{File: confFile, Message: "不能为空", Field: "Sitemap.URL"}
		case s.TagPriority > 1 || s.TagPriority < 0:
			return &FieldError{File: confFile, Message: "介于[0,1]之间的浮点数", Field: "Sitemap.TagPriority"}
		case s.PostPriority > 1 || s.PostPriority < 0:
			return &FieldError{File: confFile, Message: "介于[0,1]之间的浮点数", Field: "Sitemap.PostPriority"}
		case !isChangereq(s.TagChangefreq):
			return &FieldError{File: confFile, Message: "取值不正确", Field: "Sitemap.TagChangefreq"}
		case !isChangereq(s.PostChangefreq):
			return &FieldError{File: confFile, Message: "取值不正确", Field: "Sitemap.PostChangefreq"}
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

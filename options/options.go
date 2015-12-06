// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// options包用于处理options数据表的相关功能。
package options

import (
	"reflect"
	"strings"

	"github.com/caixw/typing/models"
	"github.com/issue9/conv"
	"github.com/issue9/orm"
	"github.com/issue9/orm/fetch"
)

// 获取评论时的返回顺序
const (
	CommentOrderUndefined = iota
	CommentOrderDesc
	CommentOrderAsc
)

// Options 对应着从options表中加载而来的数据，
// 通过struct tag来确定其对应数据库中的那条记录。
type Options struct {
	SiteName    string `options:"system,siteName"`    // 站点名称
	SecondTitle string `options:"system,secondTitle"` // 副标题
	SiteURL     string `options:"system,siteURL"`     // 网站的url
	Keywords    string `options:"system,keywords"`    // 默认页面的keywords内容
	Description string `options:"system,description"` // 默认页面的description内容
	Suffix      string `options:"system,suffix"`      // URL地址的后缀名，仅对文章有效
	//Language    string `options:"system,language"`      // 界面语言

	// 一些统计数据
	Uptime               int64  `options:"stat,uptime"`               // 上线时间
	LastUpdated          int64  `options:"stat,lastUpdated"`          // 最后更新时间
	CommentsSize         int    `options:"stat,commentsSize"`         // 评论数
	WattingCommentsSize  int    `options:"stat,wattingCommentsSize"`  // 待评论数量
	ApprovedCommentsSize int    `options:"stat,approvedCommentsSize"` // 待评论数量
	SpamCommentsSize     int    `options:"stat,spamCommentsSize"`     // 垃圾论数量
	PostsSize            int    `options:"stat,postsSize"`            // 文章数量
	PublishedPostsSize   int    `options:"stat,publishedPostsSize"`   // 已发表文章数量
	DraftPostsSize       int    `options:"stat,draftPostsSize"`       // 草稿数量
	LastLogin            int64  `options:"stat,lastLogin"`            // 最后登录时间
	LastIP               string `options:"stat,lastIP"`               // 最后登录的ip地址
	LastAgent            string `options:"stat,lastAgent"`            // 最后次登录的用户浏览器

	PageSize        int    `options:"reading,pageSize"`        // 默认每页显示的数量
	SidebarSize     int    `options:"reading,sidebarSize"`     // 侧边栏每个列表项内显示的数量
	LongDateFormat  string `options:"reading,longDateFormat"`  // 前端长日期格式
	ShortDateFormat string `options:"reading,shortDateFormat"` // 前端短日期格式
	CommentOrder    int    `options:"reading,commentOrder"`    // 评论显示的顺序

	PostsChangefreq string  `options:"feed,postsChangefreq"`
	TagsChangefreq  string  `options:"feed,tagsChangefreq"`
	PostsPriority   float64 `options:"feed,postsPriority"`
	TagsPriority    float64 `options:"feed,tagsPriority"`
	RssSize         int     `options:"feed,rssSize"` // rss和atom的数量

	Theme string `options:"themes,theme"` // 主题

	Menus string `options:"menus,menus"` // 菜单

	ScreenName string `options:"users,screenName"` // 用户昵称
	Email      string `options:"users,email"`      // 用户邮箱，可能会被显示在前端
	Password   string `options:"users,password"`   // 用户的登录密码
}

// 初始化core包。返回程序必要的变量。
func Init(db *orm.DB) (*Options, error) {
	sql := "SELECT * FROM #options"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	opt := &Options{}
	if err = opt.fromMaps(maps); err != nil {
		return nil, err
	}
	return opt, nil
}

// 将一个map数组中的数据导入到当前的options中
//  {"group":"system", "key":"pageSize", "value":"5"} ==> options.PageSize=5
func (opt *Options) fromMaps(maps []map[string]string) error {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	m := make(map[string]reflect.Value, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag.Get("options")
		m[tags] = v.Field(i)
	}

	for _, item := range maps {
		obj, found := m[item["group"]+","+item["key"]]
		if !found {
			continue
		}
		if err := conv.Value(item["value"], obj); err != nil {
			return err
		}
	}
	return nil
}

func (opt *Options) setValue(key string, val interface{}) error {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag.Get("options")
		keys := strings.Split(tags, ",")
		if keys[1] == key {
			return conv.Value(val, v.Field(i))
		}
	}
	return nil
}

// 设置options中的值，顺便更新数据库中的值。
func (opt *Options) Set(db *orm.DB, key string, val interface{}) error {
	if err := opt.setValue(key, val); err != nil {
		return err
	}

	o := &models.Option{Key: key, Value: conv.MustString(val, "")}
	if _, err := db.Update(o); err != nil {
		return err
	}

	return nil
}

// 获取指定名称的值。
func (opt *Options) Get(key string) (value interface{}, found bool) {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag.Get("options")
		index := strings.IndexByte(tags, ',')
		if tags[index+1:] == key {
			return v.Field(i).Interface(), true
		}
	}

	return nil, false
}

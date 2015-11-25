// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
	"strings"

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
	Uptime      int64  `options:"system,uptime"`      // 系统的上线时间
	//Language    string `options:"system,language"`      // 界面语言

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

// LoadOptions 从options表中加载数据，并将其转换成Options实例。
func loadOptions(db *orm.DB) (*Options, error) {
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

// 根据option的实例，更新options中某个字段，若未找到与之相对应的字段，则返回error
func (opt *Options) UpdateFromOption(o *Option) error {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag.Get("options")
		index := strings.IndexByte(tags, ',')
		if tags[index+1:] == o.Key {
			return conv.Value(o.Value, v.Field(i))
		}
	}

	return fmt.Errorf("在options实例中未找到与之[%v]相对应的字段", o.Key)
}

func (opt *Options) GetValueByKey(key string) (value interface{}, found bool) {
	// TODO 提交缓存opt的reflect.Value变量
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

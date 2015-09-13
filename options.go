// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/issue9/conv"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// 获取评论时的返回顺序
const (
	commentOrderDesc = iota
	commentOrderAsc
)

// 对应着从options表中加载而来的数据。
type options struct {
	SiteName     string `options:"system,siteName"`     // 重置的默认密码
	SiteURL      string `options:"system,siteURL"`      // 网站的url
	PageSize     int    `options:"system,pageSize"`     // 默认每页显示的数量
	Theme        string `options:"system,theme"`        // 主题
	Pretty       bool   `options:"system,pretty"`       // 格式化输出内容
	DateFormat   string `options:"system,dateFormat"`   // 前端日期格式
	CommentOrder int    `options:"system,commentOrder"` // 评论显示的顺序

	ScreenName string `options:"users,screenName"` // 用户昵称
	Email      string `options:"users,email"`      // 用户邮箱，可能会被显示在前端
	Password   string `options:"users,password"`   // 用户的登录密码
}

// 系统设置项。
type option struct {
	Key   string `orm:"name(key);len(20);pk"` // 该设置项的唯一名称
	Value string `orm:"name(value);len(-1)"`  // 该设置项的值
	Group string `orm:"name(group);len(20)"`  // 对该设置项的分组。
}

func (opt *option) Meta() string {
	return `name(options)`
}

// 将一个map数组中的数据导入到当前的options中
//  {"group":"system", "key":"pageSize", "value":"5"} ==> options.PageSize=5
func (opt *options) fromMaps(maps []map[string]string) error {
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

// 将options中的每字段转换成一个map结构。
//  options.PageSize=5 ==> {"group":"system", "key":"pageSize", "value":"5"}
func (opt *options) toMaps() ([]map[string]string, error) {
	v := reflect.ValueOf(opt)
	v = v.Elem()
	t := v.Type()
	l := t.NumField()
	maps := make([]map[string]string, 0, l)

	for i := 0; i < l; i++ {
		tags := strings.Split(t.Field(i).Tag.Get("options"), ",")
		if len(tags) != 2 {
			return nil, fmt.Errorf("len(tags)!=2 @ %v", t.Field(i).Name)
		}

		val, err := conv.String(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		maps = append(maps, map[string]string{
			"group": tags[0],
			"key":   tags[1],
			"value": val,
		})
	}

	return maps, nil
}

// 根据option的实例，更新options中某个字段，若未找到与之相对应的字段，则返回error
func (opt *options) updateFromOption(o *option) error {
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

func (opt *options) getValueByKey(key string) (value interface{}, found bool) {
	// TODO 提交缓存opt的reflect.Value变量，这用每次操作
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

// 将options表中的内容加载到opt中。
func loadOptions() (*options, error) {
	sql := "select * from #options"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}
	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	opt := &options{}
	if err = opt.fromMaps(maps); err != nil {
		return nil, err
	}
	return opt, nil
}

// @api patch /admin/api/options/{key} 修改设置项的值
// @apiParam key string 需要修改项的key
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam value string 新值
// @apiExample json
// { "value": "abcdef" }
// @apiSuccess 204 no content
func patchOption(w http.ResponseWriter, r *http.Request) {
	key, ok := paramString(w, r, "key")
	if !ok {
		return
	}

	o := &option{Key: key}
	cnt, err := db.Count(o)
	if err != nil {
		logs.Error("patchOption:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if cnt == 0 {
		renderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	if !readJSON(w, r, o) {
		return
	}

	if o.Key != key || len(o.Group) > 0 { // 提交了额外的数据内容
		renderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	// 更新数据库中的值
	if _, err := db.Update(o); err != nil {
		logs.Error("patchOption:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 更新opt变量中的值
	if err := opt.updateFromOption(o); err != nil {
		logs.Error("patchOption:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	renderJSON(w, http.StatusNoContent, nil, nil)
}

// @api get /admin/api/options/{key} 更新设置项的值，不能更新password字段。
// @apiParam key string 名称
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 ok
// @api value any 设置项的值
// @apiExample json
// {
//     "value": "20",
// }
func getOption(w http.ResponseWriter, r *http.Request) {
	key, ok := paramString(w, r, "key")
	if !ok {
		return
	}

	if key == "password" {
		renderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	val, found := opt.getValueByKey(key)
	if !found {
		renderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	renderJSON(w, http.StatusOK, map[string]interface{}{"value": val}, nil)
}

// @api get /admin/api/themes 获取所有主题列表
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
func getThemes(w http.ResponseWriter, r *http.Request) {
	// TODO
}

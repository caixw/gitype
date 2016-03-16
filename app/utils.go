// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"path"
	"strconv"

	"github.com/issue9/logs"
	"github.com/issue9/web"
)

// ParamString 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回false
func paramString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	val, found := web.ParamString(r, key)
	if found {
		return val, true
	}

	logs.Infof("app.paramString:并未在路径中找到相匹配的参数[%v]\n", val)
	w.WriteHeader(http.StatusNotFound)
	return "", false
}

// QueryInt 用于获取查询参数key的值，并将其转换成Int类型，若该值不存在返回def作为其默认值，
// 若是类型不正确，则返回一个false，并向客户端输出一个400错误。
func queryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		logs.Error("app.queryInt:", err)
		w.WriteHeader(http.StatusBadRequest)
		return 0, false
	}
	return ret, true
}

func (a *app) postURL(slug string) string {
	u := a.data.URLS
	return path.Join(u.Root, u.Post, slug+u.Suffix)
}

func (a *app) postsURL(page uint) string {
	u := a.data.URLS

	if page <= 1 {
		return u.Root
	}
	return path.Join(u.Root, u.Posts+u.Suffix) + "?page=" + strconv.Itoa(int(page))
}

func (a *app) tagURL(slug string, page uint) string {
	u := a.data.URLS

	base := path.Join(u.Root, u.Tag, slug+u.Suffix)
	if page <= 1 {
		return base
	}

	return base + "?page=" + strconv.Itoa(int(page))
}

func (a *app) tagsURL() string {
	u := a.data.URLS
	return path.Join(u.Root, u.Tags+u.Suffix)
}

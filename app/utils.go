// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"net/http"
	"strconv"

	"github.com/issue9/logs"
	"github.com/issue9/mux"
	"github.com/issue9/mux/params"
)

// 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回 false
func (a *app) paramString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	ps := mux.Params(r)
	val, err := ps.String(key)

	if err == params.ErrParamNotExists {
		a.renderError(w, http.StatusNotFound)
		return "", false
	} else if err != nil {
		logs.Error(err)
		a.renderError(w, http.StatusNotFound)
		return "", false
	} else if len(val) == 0 {
		a.renderError(w, http.StatusNotFound)
		return "", false
	}

	return val, true
}

// 获取查询参数 key 的值，并将其转换成 Int 类型，若该值不存在返回 def 作为其默认值，
// 若是类型不正确，则返回一个 false，并向客户端输出一个 400 错误。
func (a *app) queryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		logs.Error(err)
		a.renderError(w, http.StatusBadRequest)
		return 0, false
	}
	return ret, true
}

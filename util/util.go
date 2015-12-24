// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// util 包含了一些各个模块需要用到的函数。
package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"

	"github.com/issue9/web"
)

// 输出json数据到客户端。
func RenderJSON(w http.ResponseWriter, code int, v interface{}, headers map[string]string) {
	web.RenderJSON(w, code, v, headers)
}

// 从客户端读取数据到v变量中，若出错，则返回false，并向客户端输出一条错误代码。
func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) (ok bool) {
	if code := web.ReadJSON(r, v); code != http.StatusOK {
		w.WriteHeader(code)
		return false
	}

	return true
}

// ParamString 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回false
func ParamString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	if val, found := web.ParamString(r, key); found {
		return val, true
	}

	w.WriteHeader(http.StatusNotFound)
	return "", false
}

// ParamInt64 功能同ParamString，但会尝试将返回值转换成int64类型。
// 若不能找到该参数，返回false
func ParamInt64(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	if val, found := web.ParamInt64(r, key); found {
		return val, true
	}

	w.WriteHeader(http.StatusNotFound)
	return 0, false
}

// ParamID 功能同ParamInt64，但值必须大于0
func ParamID(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	if val, found := web.ParamID(r, key); found {
		return val, true
	}

	w.WriteHeader(http.StatusNotFound)
	return 0, false
}

// QueryInt 用于获取查询参数key的值，并将其转换成Int类型，若该值不存在返回def作为其默认值，
// 若是类型不正确，则返回一个false，并向客户端输出一个400错误。
func QueryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, nil, nil)
		return 0, false
	}
	return ret, true
}

// HashPassword 是一个简单的密码加密函数。
// 若需要更换密码加密算法，更改此函数即可。
func HashPassword(password string) string {
	return MD5(password)
}

// 将一段字符串转换成md5编码
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

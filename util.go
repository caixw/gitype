// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/context"
	"github.com/issue9/logs"
)

// 向客户端返回的错误信息
type ErrorResult struct {
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail,omitempty"`
}

// 将v转换成json数据并写入到w中。code为服务端返回的代码。
// 若v的值是string,[]byte,[]rune则直接转换成字符串写入w。
// 当v值为nil, "", []byte(""), []rune("")等值时，将输出表示空对象的json字符串："{}"。
// headers用于指定额外的Header信息，若传递nil，则表示没有。
func renderJSON(w http.ResponseWriter, code int, v interface{}, headers map[string]string) {
	if w == nil {
		panic("renderJSON:参数w不能为空")
	}

	var data []byte
	var err error
	switch val := v.(type) {
	case string:
		data = []byte(val)
	case []byte:
		data = val
	case []rune:
		data = []byte(string(val))
	default:
		if val == nil {
			break
		}

		if opt.Pretty {
			data, err = json.MarshalIndent(val, "", "    ")
		} else {
			data, err = json.Marshal(val)
		}
		if err != nil {
			logs.Error("renderJSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(code)
	if v == nil {
		return
	}
	if _, err = w.Write(data); err != nil {
		logs.Error("renderJSON:", err)
	}
}

// 将r中的body当作一个json格式的数据读取到v中，若出错，则直接向w输出出错内容，
// 并返回false，或是在一切正常的情况下返回true
func readJSON(w http.ResponseWriter, r *http.Request, v interface{}) (ok bool) {
	if !checkJSONMediaType(r) {
		renderJSON(w, http.StatusUnsupportedMediaType, nil, nil)
		return false
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error("readJSON:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return false
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		logs.Error("readJSON:", err)
		renderJSON(w, 422, nil, nil) // http包中并未定义422错误
		return false
	}

	return true
}

func checkJSONMediaType(r *http.Request) bool {
	if r.Method != "GET" {
		ct := r.Header.Get("Content-Type")
		if strings.Index(ct, "application/json") < 0 && strings.Index(ct, "*/*") < 0 {
			return false
		}
	}

	aceppt := r.Header.Get("Accept")
	return strings.Index(aceppt, "application/json") >= 0 || strings.Index(aceppt, "*/*") >= 0
}

// 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回false
func paramString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	m, found := context.Get(r).Get("params")
	if !found {
		logs.Error("paramString:在context中找不到params参数:params")
		renderJSON(w, http.StatusGone, nil, nil)
		return "", false
	}

	params := m.(map[string]string)
	v, found := params[key]
	if !found {
		logs.Error("paramString:在context.params中找不到指定参数:", key)
		renderJSON(w, http.StatusGone, nil, nil)
		return "", false
	}

	return v, true
}

// 获取路径匹配中的参数，并以int64的格式返回。
// 若不能找到该参数，返回false
func paramInt64(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	v, ok := paramString(w, r, key)
	if !ok {
		return 0, false
	}

	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logs.Error("paramInt64:", err)
		renderJSON(w, http.StatusGone, nil, nil)
		return 0, false
	}

	return num, true
}

func paramID(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	num, ok := paramInt64(w, r, key)
	if !ok {
		return 0, false
	}

	if num <= 0 {
		logs.Error("paramID:用户指定了一个小于0的id值:", num)
		renderJSON(w, http.StatusGone, nil, nil)
		return 0, false
	}

	return num, true
}

// 获取查询参数key的值，若该值不存在返回def作为其默认值，
// 若是类型不正确，则返回一个false，并向客户端输出一个400错误。
func queryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, nil, nil)
		return 0, false
	}
	return ret, true
}

// 简单的密码加密函数
func hashPassword(password string) string {
	m := md5.New()
	m.Write([]byte(password))
	return hex.EncodeToString(m.Sum(nil))
}

// 将换行符转换成<br />标签
func nl2br(str string) string {
	return strings.Replace(str, "\n", "<br />", -1)
}

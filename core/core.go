// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/issue9/context"
	"github.com/issue9/logs"
	"github.com/issue9/orm"
	"github.com/issue9/orm/dialect"
	"github.com/issue9/orm/forward"
	"github.com/issue9/web"
)

const (
	Version = "0.4.30.151122" // 程序版本号

	// 两个配置文件路径
	ConfigPath    = "./config/app.json"
	LogConfigPath = "./config/logs.xml"
)

var (
	Cfg *Config
	Opt *Options
	DB  *orm.DB
)

// 初始化core包。返回程序必要的变量。
func Init() (err error) {
	Cfg, err = LoadConfig(ConfigPath)
	if err != nil {
		return
	}

	DB, err = InitDB(Cfg)
	if err != nil {
		return
	}

	if err = logs.InitFromXMLFile(LogConfigPath); err != nil {
		return
	}

	Opt, err = loadOptions(DB)
	return
}

func Run() {
	web.Run(Cfg.Core)
}

func Close() {
	DB.Close()
}

// RenderJSON 用于将v转换成json数据并写入到w中。code为服务端返回的代码。
// 若v的值是string,[]byte,[]rune则直接转换成字符串写入w。
// 当v为nil时，不输出任何内容，若需要输出一个空对象，请使用"{}"字符串。
// headers用于指定额外的Header信息，若传递nil，则表示没有。
func RenderJSON(w http.ResponseWriter, code int, v interface{}, headers map[string]string) {
	if w == nil {
		panic("RenderJSON:参数w不能为空")
	}

	if v == nil {
		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		for k, v := range headers {
			w.Header().Add(k, v)
		}
		w.WriteHeader(code)
		return
	}

	var data []byte
	switch val := v.(type) {
	case string:
		data = []byte(val)
	case []byte:
		data = val
	case []rune:
		data = []byte(string(val))
	default:
		var err error
		data, err = json.Marshal(val)
		if err != nil {
			logs.Error("RenderJSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		logs.Error("RenderJSON:", err)
	}
}

// ReadJSON 用于将r中的body当作一个json格式的数据读取到v中，若出错，则直接向w输出出错内容，
// 并返回false，或是在一切正常的情况下返回true
func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) (ok bool) {
	if !checkJSONMediaType(r) {
		RenderJSON(w, http.StatusUnsupportedMediaType, nil, nil)
		return false
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error("readJSON:", err)
		RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return false
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		logs.Error("readJSON:", err)
		RenderJSON(w, 422, nil, nil) // http包中并未定义422错误
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

// ParamString 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回false
func ParamString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	m, found := context.Get(r).Get("params")
	if !found {
		logs.Error("ParamString:在context中找不到params参数")
		RenderJSON(w, http.StatusGone, nil, nil)
		return "", false
	}

	params := m.(map[string]string)
	v, found := params[key]
	if !found {
		logs.Error("ParamString:在context.params中找不到指定参数:", key)
		RenderJSON(w, http.StatusGone, nil, nil)
		return "", false
	}

	return v, true
}

// ParamInt64 功能同ParamString，但会尝试将返回值转换成int64类型。
// 若不能找到该参数，返回false
func ParamInt64(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	v, ok := ParamString(w, r, key)
	if !ok {
		return 0, false
	}

	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logs.Error("ParamInt64:", err)
		RenderJSON(w, http.StatusGone, nil, nil)
		return 0, false
	}

	return num, true
}

// ParamID 功能同ParamInt64，但值必须大于0
func ParamID(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	num, ok := ParamInt64(w, r, key)
	if !ok {
		return 0, false
	}

	if num <= 0 {
		logs.Trace("ParamID:用户指定了一个小于0的id值:", num)
		RenderJSON(w, http.StatusNotFound, nil, nil)
		return 0, false
	}

	return num, true
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
// 若需要更换密码加密算法，理发此函数即可。
func HashPassword(password string) string {
	return MD5(password)
}

// 将一段字符串转换成md5编码
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 从一个Config实例中初始一个orm.DB实例。
func InitDB(cfg *Config) (*orm.DB, error) {
	var d forward.Dialect
	switch cfg.DBDriver {
	case "sqlite3":
		d = dialect.Sqlite3()
	case "mysql":
		d = dialect.Mysql()
	case "postgres":
		d = dialect.Postgres()
	default:
		return nil, errors.New("不能理解的dbDriver值：" + cfg.DBDriver)
	}

	return orm.NewDB(cfg.DBDriver, cfg.DBDSN, cfg.DBPrefix, d)
}

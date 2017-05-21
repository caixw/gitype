// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"path"
	"strconv"

	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
)

// 获取路径匹配中的参数，并以字符串的格式返回。
// 若不能找到该参数，返回 false
func (c *Client) paramString(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	ps := mux.GetParams(r)
	val, err := ps.String(key)

	if err == mux.ErrParamNotExists {
		c.renderError(w, http.StatusBadRequest)
		return "", false
	} else if err != nil {
		logs.Error(err)
		c.renderError(w, http.StatusBadRequest)
		return "", false
	} else if len(val) == 0 {
		c.renderError(w, http.StatusBadRequest)
		return "", false
	}

	return val, true
}

// 获取查询参数 key 的值，并将其转换成 Int 类型，若该值不存在返回 def 作为其默认值，
// 若是类型不正确，则返回一个 false，并向客户端输出一个 400 错误。
func (c *Client) queryInt(w http.ResponseWriter, r *http.Request, key string, def int) (int, bool) {
	val := r.FormValue(key)
	if len(val) == 0 {
		return def, true
	}

	ret, err := strconv.Atoi(val)
	if err != nil {
		logs.Error(err)
		c.renderError(w, http.StatusBadRequest)
		return 0, false
	}
	return ret, true
}

func (c *Client) postURL(slug string) string {
	return path.Join(vars.Post, slug+vars.Suffix)
}

func (c *Client) postsURL(page int) string {
	if page <= 1 {
		return "/"
	}
	return vars.Posts + vars.Suffix + "?page=" + strconv.Itoa(page)
}

func (c *Client) tagURL(slug string, page int) string {
	url := path.Join(vars.Tag, slug+vars.Suffix)
	if page <= 1 {
		return url
	}

	return url + "?page=" + strconv.Itoa(page)
}

func (c *Client) searchURL(q string, page int) string {
	url := vars.Search + vars.Suffix
	if len(q) > 0 {
		url += "?q=" + q
	}

	if page > 1 {
		if len(q) > 0 {
			url += "&"
		} else {
			url += "?"
		}
		url += "page=" + strconv.Itoa(page)
	}

	return url
}

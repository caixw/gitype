// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/issue9/logs"
)

var (
	tplDir string
	themes map[string]*theme
	tpl    *template.Template
)

type theme struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	Screenshot  string  `json:"screenshot"`
	Author      *author `json:"author"`
}

type author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// 初始化主题相关的操作。
func initThemes(themeDir string) error {
	tplDir = themeDir + string(os.PathSeparator)
	fs, err := ioutil.ReadDir(tplDir)
	if err != nil {
		return err
	}

	if len(fs) == 0 {
	}

	themes = make(map[string]*theme, len(fs))

	for _, file := range fs {
		if !file.IsDir() {
			continue
		}

		path := themeDir + "/" + file.Name() + "/theme.json"
		t, err := loadThemeFile(path)
		if err != nil {
			return err
		}

		themes[file.Name()] = t
	}

	return nil
}

func loadThemeFile(path string) (*theme, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := &theme{}
	if err = json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	return t, nil
}

func loadCurrTheme() (err error) {
	tpl, err = template.ParseGlob(tplDir + opt.Theme + "/*.html")
	return
}

// @api get /admin/api/themes 获取所有主题列表
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam themes array 所有主题列表
func adminGetThemes(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, http.StatusOK, map[string]interface{}{"themes": themes}, nil)
}

// @api patch /admin/api/themes/current 更改当前的主题
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam value string 新值
//
// @apiSuccess 200 OK
func adminPostTheme(w http.ResponseWriter, r *http.Request) {
	o := &option{Key: "theme"}
	if !readJSON(w, r, o) {
		return
	}

	if o.Key != "theme" || len(o.Group) > 0 { // 提交了额外的数据内容
		renderJSON(w, http.StatusBadRequest, nil, nil)
		return
	}

	if err := patchOption(o); err != nil {
		logs.Error("adminPostTheme:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := loadCurrTheme(); err != nil {
		logs.Error("adminPostTheme:", err)
		renderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	renderJSON(w, http.StatusNoContent, nil, nil)
}

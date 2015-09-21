// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

// 主题管理
type Themes struct {
	dir       string             // 主题目录
	urlPrefix string             // 主题的URL前缀
	tpl       *template.Template // 当前使用的模板
	themes    map[string]*Theme  // 所有的主题列表

}

// Theme 用于描述主题的相关信息，一般从主题目录下的theme.json获取。
type Theme struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	Screenshot  string  `json:"screenshot"`
	Author      *Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// 从主题根目录加载所有的主题内容，并初始所有的主题下静态文件的路由。
// defaultTheme 为默认的主题。
func LoadThemes(cfg *Config, defaultTheme string) (*Themes, error) {
	dir := cfg.ThemeDir + string(os.PathSeparator)
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(fs) == 0 {
		return nil, errors.New("不存在任何主题目录")
	}

	themes := &Themes{
		dir:       dir,
		urlPrefix: cfg.ThemeURLPrefix,
		themes:    make(map[string]*Theme, len(fs)),
	}

	p := cfg.ThemeURLPrefix + "/"
	for _, file := range fs {
		if !file.IsDir() {
			continue
		}
		name := file.Name()
		themePath := dir + name + string(os.PathSeparator)

		path := themePath + "theme.json"
		t, err := loadThemeFile(path)
		if err != nil {
			return nil, err
		}

		t.Screenshot = p + name + "/" + t.Screenshot
		themes.themes[name] = t
		cfg.Core.Static[p+name] = themePath + "public/"
	}

	return themes, themes.Switch(defaultTheme)
}

func loadThemeFile(path string) (*Theme, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := &Theme{}
	if err = json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	return t, nil
}

// 返回所有的主题列表
func (t *Themes) Themes() map[string]*Theme {
	return t.themes
}

// 切换主题
func (t *Themes) Switch(id string) (err error) {
	t.tpl, err = template.ParseGlob(t.dir + id + "/*.html")
	return
}

// 输出指定模板
func (t *Themes) Render(w http.ResponseWriter, name string, data interface{}) error {
	return t.tpl.ExecuteTemplate(w, name, data)
}

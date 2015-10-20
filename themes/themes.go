// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/caixw/typing/core"
	"github.com/issue9/logs"
)

var (
	cfg       *core.Config
	tpl       *template.Template // 当前使用的模板
	themesMap map[string]*Theme  // 所有的主题列表
)

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
func Init(config *core.Config, defaultTheme string) error {
	cfg = config

	fs, err := ioutil.ReadDir(cfg.ThemeDir)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return errors.New("不存在任何主题目录")
	}
	themesMap = make(map[string]*Theme, len(fs))

	p := cfg.ThemeURLPrefix + "/"
	for _, file := range fs {
		if !file.IsDir() {
			continue
		}

		name := file.Name()
		themePath := cfg.ThemeDir + name + string(os.PathSeparator)

		theme, err := loadThemeFile(themePath + "theme.json")
		if err != nil {
			return err
		}

		theme.Screenshot = p + name + "/" + theme.Screenshot
		themesMap[name] = theme
		cfg.Core.Static[p+name] = themePath + "public/"
	}

	return Switch(defaultTheme)
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
func Themes() map[string]*Theme {
	return themesMap
}

// 切换主题
func Switch(id string) (err error) {
	tpl, err = template.ParseGlob(cfg.ThemeDir + id + "/*.html")
	return
}

// 输出指定模板
func Render(w http.ResponseWriter, name string, data interface{}) {
	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		logs.Error("core.Render:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
}

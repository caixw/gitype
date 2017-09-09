// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
	"github.com/issue9/mux"
	"github.com/issue9/utils"
)

// 模板的扩展名，在主题目录下，以下扩展名的文件，不会被展示
var ignoreThemeFileExts = []string{
	vars.TemplateExtension,
	".yaml",
	".yml",
}

func isIgnoreThemeFile(file string) bool {
	ext := filepath.Ext(file)

	for _, v := range ignoreThemeFileExts {
		if ext == v {
			return true
		}
	}

	return false
}

// 资源内容
// /posts/{path}
func (client *Client) getAsset(w http.ResponseWriter, r *http.Request) {
	path, err := mux.Params(r).String("path")
	if err != nil {
		logs.Error(err)
		client.getRaws(w, r)
		return
	}

	// 不展示模板文件，查看 raws 中是否有同名文件
	name := filepath.Base(path)
	if name == vars.PostMetaFilename || name == vars.PostContentFilename {
		client.getRaws(w, r)
		return
	}

	filename := filepath.Join(client.path.PostsDir, path)

	if !utils.FileExists(filename) {
		client.getRaws(w, r)
		return
	}

	stat, err := os.Stat(filename)
	if err != nil {
		logs.Error(err)
		client.renderError(w, http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		client.getRaws(w, r)
		return
	}

	http.ServeFile(w, r, filename)
}

// 主题文件
// /themes/...
func (client *Client) getThemes(w http.ResponseWriter, r *http.Request) {
	if isIgnoreThemeFile(r.URL.Path) { // 不展示模板文件，查看 raws 中是否有同名文件
		client.getRaws(w, r)
		return
	}

	path, err := mux.Params(r).String("path")
	if err != nil {
		logs.Error(err)
		client.getRaws(w, r)
		return
	}

	filename := filepath.Join(client.path.ThemesDir, path)

	if !utils.FileExists(filename) {
		client.getRaws(w, r)
		return
	}

	stat, err := os.Stat(filename)
	if err != nil {
		logs.Error(err)
		client.renderError(w, http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		client.getRaws(w, r)
		return
	}

	http.ServeFile(w, r, filename)
}

// /...
func (client *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		client.getPosts(w, r)
		return
	}

	if !utils.FileExists(filepath.Join(client.path.RawsDir, r.URL.Path)) {
		client.renderError(w, http.StatusNotFound)
		return
	}

	prefix := "/"
	root := http.Dir(client.path.RawsDir)
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)
}

// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"path/filepath"

	"github.com/caixw/typing/vars"
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

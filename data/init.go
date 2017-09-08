// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"os"

	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// Init 在 path 下初始化基本的数据
func Init(path *vars.Path) error {
	if !utils.FileExists(path.DataDir) {
		if err := os.Mkdir(path.DataDir, os.ModePerm); err != nil {
			return err
		}
	}

	if !utils.FileExists(path.MetaDir) {
		if err := os.Mkdir(path.MetaDir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

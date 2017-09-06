// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package init

import (
	"os"

	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// Init 指行初始化命令
func Init(path *vars.Path) error {
	if err := initConfDir(path); err != nil {
		return err
	}

	return initDataDir(path)
}

// 初始化 conf 目录下的数据
func initConfDir(path *vars.Path) error {
	if !utils.FileExists(path.ConfDir) {
		if err := os.Mkdir(path.ConfDir, os.ModePerm); err != nil {
			return err
		}
	}

	// TODO

	return nil
}

func initDataDir(path *vars.Path) error {
	if !utils.FileExists(path.DataDir) {
		if err := os.Mkdir(path.DataDir, os.ModePerm); err != nil {
			return err
		}
	}

	// TODO

	return nil
}

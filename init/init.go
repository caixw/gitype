// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package init

import (
	"fmt"
	"os"

	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
)

// Init 执行初始化命令
func Init(path *vars.Path) error {
	if err := initConfDir(path); err != nil {
		return err
	}

	if err := initDataDir(path); err != nil {
		return err
	}

	_, err := fmt.Fprintf(vars.CMDOutput, "操作成功，你现在可以在 %s 中修改具体的参数配置！", path.Root)
	return err
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

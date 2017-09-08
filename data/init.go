// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"os"
	"time"

	"github.com/caixw/typing/vars"
	"github.com/issue9/utils"
	yaml "gopkg.in/yaml.v2"
)

var defaultConfig = &config{
	Title:           "Title",
	Language:        "zh-cnm-Hans",
	Subtitle:        "subtitle",
	URL:             "https://caixw.io",
	Keywords:        vars.AppName,
	PageSize:        20,
	LongDateFormat:  "2006-01-02 15:04:05",
	ShortDateFormat: "2006-01-02",
	Author: &Author{
		Name: vars.AppName,
		URL:  vars.URL,
	},
	License: &Link{
		Title: "署名 4.0 国际 (CC BY 4.0)",
		URL:   "https://creativecommons.org/licenses/by/4.0/deed.zh",
	},

	Theme:        "default",
	UptimeFormat: time.Now().Format(vars.DateFormat),
	Archive: &archiveConfig{
		Type:   archiveTypeYear,
		Format: "2006 年",
	},
}

// Init 在 path 下初始化基本的数据
func Init(path *vars.Path) error {
	fmt.Println(path.DataDir)
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

	// data/meta/config.yaml
	file, err := os.Create(path.MetaConfigFile)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	if _, err := file.Write(bs); err != nil {
		return err
	}

	return nil
}

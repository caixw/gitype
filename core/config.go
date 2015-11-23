// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/issue9/mux"
	"github.com/issue9/web"
)

// Config 表示程序级别的配置，修改这些配置需要重启程序才能启作用，
// 比如数据库初始化信息，路由项设置等。
type Config struct {
	Core *web.Config `json:"core"`

	Debug   bool   `json:"debug"`   // 是否处于调试模式
	TempDir string `json:"tempDir"` // 临时文件所在的目录，该目录下的文件被删除不会影响程序整体运行。

	// 数据库相关配置
	DBDSN    string `json:"dbDSN"`    // 数据库dsn
	DBPrefix string `json:"dbPrefix"` // 数据表名前缀
	DBDriver string `json:"dbDriver"` // 数据库类型，可以是mysql, sqlite3, postgresql

	// 接口相关配置
	FrontAPIPrefix string `json:"frontAPIPrefix"` // 前端api地址前缀
	AdminAPIPrefix string `json:"adminAPIPrefix"` // 后台api地址前缀

	// 主题相关配置
	ThemeURLPrefix string `json:"themeURLPrefix"` // 各主题公开文件的根URL
	ThemeDir       string `json:"themeDir"`       // 主题文件所在的目录

	// 上传文件相关配置
	UploadDir       string `json:"uploadDir"`       // 上传文件所在的目录
	UploadDirFormat string `json:"uploadDirFormat"` // 上传文件子路径的格式，只能以时间为格式
	UploadSize      int64  `json:"uploadSize"`      // 上传文件的最大尺寸
	UploadExts      string `json:"uploadExts"`      // 允许的上传文件扩展名，eg: .txt;.png,;.pdf
	UploadURLPrefix string `json:"uploadURLPrefix"` // 上传文件的地址前缀
}

// LoadConfig 用于加载path的内容，并尝试将其转换成Config实例。
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// 检测AdminAPIPrefix是否符合要求
	if len(cfg.AdminAPIPrefix) == 0 {
		return nil, errors.New("必须指定adminApiPrefix值")
	}
	if strings.HasSuffix(cfg.AdminAPIPrefix, "/") {
		return nil, errors.New("adminAPIPrefix不能以/符号结尾")
	}

	// 检测FrontAPIPrefix是否符合要求
	if len(cfg.FrontAPIPrefix) == 0 {
		return nil, errors.New("必须指定frontApiPrefix值")
	}
	if strings.HasSuffix(cfg.FrontAPIPrefix, "/") {
		return nil, errors.New("frontAPIPrefix不能以/符号结尾")
	}

	// 检测ThemeURLPrefix是否符合要求
	if len(cfg.ThemeURLPrefix) == 0 {
		return nil, errors.New("必须指定themeURLPrefix值")
	}
	if strings.HasSuffix(cfg.ThemeURLPrefix, "/") {
		return nil, errors.New("themeURLPrefix不能以/符号结尾")
	}

	if len(cfg.ThemeDir) == 0 {
		return nil, errors.New("themeDir未指定")
	}
	if !strings.HasSuffix(cfg.ThemeDir, "/") && !strings.HasSuffix(cfg.ThemeDir, string(os.PathSeparator)) {
		return nil, errors.New("themeDir只能以/结尾")
	}

	if len(cfg.TempDir) == 0 {
		return nil, errors.New("tempDir未指定")
	}
	if !strings.HasSuffix(cfg.TempDir, "/") && !strings.HasSuffix(cfg.TempDir, string(os.PathSeparator)) {
		return nil, errors.New("tempDir只能以/结尾")
	}
	if len(cfg.DBDSN) == 0 {
		return nil, errors.New("app.json中未指定dbDSN")
	}
	if len(cfg.DBDriver) == 0 {
		return nil, errors.New("app.json中未指定dbDriver")
	}

	// upload
	if len(cfg.UploadDir) == 0 {
		return nil, errors.New("uploadDir未指定")
	}
	if len(cfg.UploadDirFormat) == 0 {
		cfg.UploadDirFormat = "2006/01/02/"
	}
	if strings.Index(cfg.UploadDirFormat, "..") >= 0 {
		return nil, errors.New("uploadDirFormat不能包含..字符")
	}
	if !strings.HasSuffix(cfg.UploadDir, "/") && !strings.HasSuffix(cfg.UploadDir, string(os.PathSeparator)) {
		return nil, errors.New("uploadDir只能以/结尾")
	}
	if cfg.UploadSize < 0 {
		return nil, errors.New("uploadSize必须大于等于0")
	}
	if strings.IndexByte(cfg.UploadExts, '/') >= 0 || strings.IndexByte(cfg.UploadExts, os.PathSeparator) >= 0 {
		return nil, errors.New("uploadExts包含非法的字符")
	}
	if len(cfg.UploadURLPrefix) == 0 {
		return nil, errors.New("必须指定uploadURLPrefix值")
	}
	if strings.HasSuffix(cfg.UploadURLPrefix, "/") {
		return nil, errors.New("uploadURLPrefix不能以/符号结尾")
	}

	if cfg.Debug { // 调试状态，输出详细信息
		cfg.Core.ErrHandler = mux.PrintDebug
	}

	return cfg, nil
}

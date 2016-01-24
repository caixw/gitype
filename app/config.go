// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/issue9/handlers"
	"github.com/issue9/utils"
	"github.com/issue9/web"
)

// Config 表示程序级别的配置，修改这些配置需要重启程序才能启作用，
// 比如数据库初始化信息，路由项设置等。
type Config struct {
	Core *web.Config `json:"core"`

	Debug bool `json:"debug"` // 是否处于调试模式

	// 后台相关的设置
	AdminURLPrefix string `json:"adminURLPrefix"` // 后台地址入口
	AdminDir       string `json:"adminDir"`       // 后台静态文件对应的目录
	Salt           string `json:"salt"`           // 密码加盐值，一量确认，不能修改

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

	RootDir string `json:"rootDir"` // 根地址下文件对应的目录

	// 上传文件相关配置
	UploadDir       string `json:"uploadDir"`       // 上传文件所在的目录
	UploadDirFormat string `json:"uploadDirFormat"` // 上传文件子路径的格式，只能以时间为格式
	UploadSize      int64  `json:"uploadSize"`      // 上传文件的最大尺寸
	UploadExts      string `json:"uploadExts"`      // 允许的上传文件扩展名，eg: .txt;.png,;.pdf
	UploadURLPrefix string `json:"uploadURLPrefix"` // 上传文件的地址前缀
}

// 获取Config实例
func GetConfig() *Config {
	return config
}

// 返回一个加盐值的密码。
func Password(password string) string {
	return utils.MD5(utils.MD5(password) + config.Salt)
}

// 检测配置项的URL，是否符合要求。
func checkConfigURL(url, field string) error {
	if len(url) == 0 {
		return fmt.Errorf("字段[%v]不能为空", field)
	}

	if strings.HasSuffix(url, "/") {
		return fmt.Errorf("字段[%v]不能以/符号结尾", field)
	}

	return nil
}

// 检测配置项的路径值，是否符合要求。
func checkConfigDir(dir, field string) error {
	if len(dir) == 0 {
		return fmt.Errorf("字段[%v]不能为空", field)
	}

	if !strings.HasSuffix(dir, "/") && !strings.HasSuffix(dir, string(os.PathSeparator)) {
		return fmt.Errorf("字段[%v]只能以路径分隔符(/或\\)作结尾", field)
	}

	return nil
}

// 加载path的内容，并尝试将其转换成Config实例。
func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	if err = checkConfigURL(cfg.AdminURLPrefix, "adminURLPrefix"); err != nil {
		return nil, err
	}

	if err = checkConfigURL(cfg.AdminAPIPrefix, "adminAPIPrefix"); err != nil {
		return nil, err
	}
	if err = checkConfigURL(cfg.FrontAPIPrefix, "frontApiPrefix"); err != nil {
		return nil, err
	}

	// theme
	if err = checkConfigURL(cfg.ThemeURLPrefix, "themeURLPrefix"); err != nil {
		return nil, err
	}
	if err = checkConfigDir(cfg.ThemeDir, "themeDir"); err != nil {
		return nil, err
	}

	if err = checkConfigDir(cfg.RootDir, "rootDir"); err != nil {
		return nil, err
	}

	// DB
	if len(cfg.DBDSN) == 0 {
		return nil, errors.New("app.json中未指定dbDSN")
	}
	if len(cfg.DBDriver) == 0 {
		return nil, errors.New("app.json中未指定dbDriver")
	}

	// upload
	if err = checkConfigDir(cfg.UploadDir, "uploadDir"); err != nil {
		return nil, err
	}
	if len(cfg.UploadDirFormat) == 0 {
		cfg.UploadDirFormat = "2006/01/02/"
	}
	if strings.Index(cfg.UploadDirFormat, "..") >= 0 {
		return nil, errors.New("uploadDirFormat不能包含..字符")
	}
	if cfg.UploadSize < 0 {
		return nil, errors.New("uploadSize必须大于等于0")
	}
	if strings.IndexByte(cfg.UploadExts, '/') >= 0 || strings.IndexByte(cfg.UploadExts, os.PathSeparator) >= 0 {
		return nil, errors.New("uploadExts包含非法的字符")
	}
	if err = checkConfigURL(cfg.UploadURLPrefix, "uploadURLPrefix"); err != nil {
		return nil, err
	}

	if cfg.Debug { // 调试状态，输出详细信息
		cfg.Core.ErrHandler = handlers.PrintDebug
	}

	return cfg, nil
}

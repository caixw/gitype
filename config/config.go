// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package config 全局的配置内容。
package config

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/helper"
	"github.com/caixw/typing/vars"
	"github.com/issue9/is"
	"github.com/issue9/utils"
)

// 两个默认端口的定义
const (
	HTTPPort  = ":80"
	HTTPSPort = ":443"
)

// 对 Config.HTTPState 可选值的定义
const (
	HTTPStateDefault  = "default"
	HTTPStateDisable  = "disable"
	HTTPStateRedirect = "redirect"
)

const (
	cookieMaxAge = 3600
)

// Config 程序的全局配置内容
type Config struct {
	// 是否启用 HTTPS 模式。如果启用了，则需要正确设置以下几个值：
	// HTTPState、CertFile、KeyFile
	HTTPS bool `yaml:"https,omitempty"`

	// 当启用 HTTPS 且端口不为 80 时，对 80 端口的处理方式。
	// disable 表示禁用 80 端口；
	// default 默认方式，即和当前的处理方式相同；
	// redirect 跳转到当前端口；
	HTTPState string `yaml:"httpState,omitempty"`

	CertFile string `yaml:"certFile,omitempty"`

	KeyFile string `yaml:"keyFile,omitempty"`

	// 监听的端口，需要带前缀冒号(:)，不指定时，
	// 根据 HTTPS 的值，默认为 :80 或是 :443
	Port string `yaml:"port,omitempty"`

	// CookieMaxAge cookie 的生存期，单位：秒
	CookieMaxAge int `yaml:"cookieMaxAge,omitempty"`

	// 绑定的域名，若指定了该值，则只能通过这些域名才能访问网站。
	// 为空表示不作限制。
	Domains []string `yaml:"domains,omitempty"`

	// Headers 用于指定一些返回给客户端的固定报头内容。
	// 其中键名表示报头名称，键值表示报头的值。
	Headers map[string]string `yaml:"headers,omitempty"`

	Webhook *Webhook `yaml:"webhook"`
}

// Webhook 与 Webhooks 相关的配置内容
type Webhook struct {
	URL       string        `yaml:"url"`              // webhooks 接收地址
	Frequency time.Duration `yaml:"frequency"`        // webhooks 的最小更新频率
	Method    string        `yaml:"method,omitempty"` // webhooks 的请求方式，默认为 POST
	RepoURL   string        `yaml:"repoURL"`          // 远程仓库的地址
}

// Load 加载配置内容
func Load(path *vars.Path) (*Config, error) {
	conf := &Config{}
	if err := helper.LoadYAMLFile(path.AppConfigFile, conf); err != nil {
		return nil, err
	}

	if err := conf.sanitize(); err != nil {
		err.File = path.AppConfigFile
		return nil, err
	}

	return conf, nil
}

func (w *Webhook) sanitize() *helper.FieldError {
	if len(w.Method) == 0 {
		w.Method = http.MethodPost
	}

	switch {
	case len(w.URL) == 0 || w.URL[0] != '/':
		return &helper.FieldError{Field: "webhook.url", Message: "不能为空且只能以 / 开头"}
	case w.Frequency < 0:
		return &helper.FieldError{Field: "webhook.frequency", Message: "不能小于 0"}
	case len(w.RepoURL) == 0:
		return &helper.FieldError{Field: "webhook.repoURL", Message: "不能为空"}
	}

	return nil
}

func (conf *Config) sanitize() *helper.FieldError {
	if len(conf.Port) == 0 {
		if conf.HTTPS {
			conf.Port = HTTPSPort
		} else {
			conf.Port = HTTPPort
		}
	}

	if conf.HTTPS {
		if len(conf.HTTPState) == 0 {
			conf.HTTPState = HTTPStateDefault
		}

		switch {
		case conf.HTTPState != HTTPStateDefault &&
			conf.HTTPState != HTTPStateDisable &&
			conf.HTTPState != HTTPStateRedirect:
			return &helper.FieldError{Field: "httpState", Message: "无效的取值"}
		case conf.HTTPState != HTTPStateDisable && conf.Port == HTTPPort:
			return &helper.FieldError{Field: "port", Message: "80 端口已经被被监听"}
		case !utils.FileExists(conf.CertFile):
			return &helper.FieldError{Field: "certFile", Message: "不能为空"}
		case !utils.FileExists(conf.KeyFile):
			return &helper.FieldError{Field: "keyFile", Message: "不能为空"}
		}
	}

	if conf.CookieMaxAge < 0 {
		return &helper.FieldError{Field: "cookieMaxAge", Message: "必须大于 0"}
	}
	if conf.CookieMaxAge == 0 {
		conf.CookieMaxAge = cookieMaxAge
	}

	if len(conf.Domains) > 0 {
		for index, domain := range conf.Domains {
			if !is.URL(domain) {
				return &helper.FieldError{Field: "domains[" + strconv.Itoa(index) + "]", Message: "无效的 URL"}
			}
		}
	}

	return conf.Webhook.sanitize()
}

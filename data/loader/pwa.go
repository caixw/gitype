// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package loader

import (
	"github.com/issue9/web"

	"github.com/caixw/gitype/helper"
)

// PWA 表示 PWA 中的相关配置
type PWA struct {
	// 指定 sw.js 的路径，为空表示不启用
	ServiceWorker string    `yaml:"serviceWorker,omitempty"`
	Manifest      *Manifest `yaml:"manifest,omitempty"`
}

// Manifest 表示 PWA 中 manifest 的相关配置
type Manifest struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type,omitempty"`

	Lang        string  `yaml:"lang"`
	Name        string  `yaml:"name"`
	ShortName   string  `yaml:"shortName"`
	StartURL    string  `yaml:"startURL,omitempty"`
	Display     string  `yaml:"display,omitempty"`
	Description string  `yaml:"description,omitempty"`
	Dir         string  `yaml:"dir,omitempty"`
	Orientation string  `yaml:"orientation,omitempty"`
	Scope       string  `yaml:"scope,omitempty"`
	ThemeColor  string  `yaml:"themeColor,omitempty"`
	Background  string  `yaml:"backgroundColor,omitempty"`
	Icons       []*Icon `yaml:"icons"`
}

func (pwa *PWA) sanitize(conf *Config) *helper.FieldError {
	if pwa.ServiceWorker != "" {
		if pwa.ServiceWorker[0] != '/' || len(pwa.ServiceWorker) == 1 {
			return &helper.FieldError{Message: "只能以 / 开头，且必须有内容", Field: "pwa.serviceWorker"}
		}
	}

	if pwa.Manifest != nil {
		return pwa.Manifest.sanitize(conf)
	}
	return nil
}

func (m *Manifest) sanitize(conf *Config) *helper.FieldError {
	if m.URL == "" {
		return &helper.FieldError{Message: "不能为空", Field: "pwa.manifest.url"}
	}

	if m.Type == "" {
		m.Type = contentManifest
	}

	if m.Lang == "" {
		m.Lang = conf.Language
	}

	if m.Name == "" {
		m.Name = conf.Title
	}

	if m.ShortName == "" {
		m.ShortName = conf.Subtitle
	}

	if m.StartURL == "" {
		m.StartURL = web.URL("")
	}

	if m.Display == "" {
		m.Display = "browser"
	} else if !inStrings(m.Display, pwaDisplays) {
		return &helper.FieldError{Message: "取值不正确", Field: "pwa.manifest.display"}
	}

	if m.Dir == "" {
		m.Dir = "auto"
	} else if !inStrings(m.Dir, pwaDirs) {
		return &helper.FieldError{Message: "取值不正确", Field: "pwa.manifest.dir"}
	}

	if m.Orientation != "" && !inStrings(m.Orientation, pwaOrientations) {
		return &helper.FieldError{Message: "取值不正确", Field: "pwa.manifest.orientation"}
	}

	if len(m.Icons) == 0 { // nil 或是 len(m.Icons) == 0
		m.Icons = []*Icon{conf.Icon}
	}

	return nil
}

var pwaDisplays = []string{
	"fullscreen",
	"standalone",
	"minimal-ul",
	"browser",
}

var pwaOrientations = []string{
	"any",
	"natural",
	"landscape",
	"landscape-primary",
	"landscape-secondary",
	"portrait",
	"portrait-primary",
	"portrait-secondary",
}

var pwaDirs = []string{
	"rtl",
	"ltr",
	"auto",
}

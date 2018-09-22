// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"encoding/json"

	"github.com/caixw/gitype/data/loader"
)

// Manifest 表示 PWA 中的 manifest.json 文件
type Manifest struct {
	Lang        string  `json:"lang"`
	Name        string  `json:"name"`
	ShortName   string  `json:"short_name"`
	StartURL    string  `json:"start_url,omitempty"`
	Display     string  `json:"display,omitempty"`
	Description string  `json:"description,omitempty"`
	Dir         string  `json:"dir,omitempty"`
	Orientation string  `json:"orientation,omitempty"`
	Scope       string  `json:"scope,omitempty"`
	ThemeColor  string  `json:"theme_color,omitempty"`
	Background  string  `json:"background_color,omitempty"`
	Icons       []*icon `json:"icons"`
}

type icon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

func (d *Data) buildPWA(conf *loader.Config) error {
	if conf.PWA == nil { // 不需要生成 pwa
		return nil
	}

	m := &Manifest{}
	m.fromLoader(conf.PWA)

	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}

	d.PWA = &Feed{
		URL:     conf.PWA.URL,
		Type:    conf.PWA.Type,
		Content: bs,
	}

	return nil
}

func (m *Manifest) fromLoader(conf *loader.Manifest) {
	m.Lang = conf.Lang
	m.Name = conf.Name
	m.ShortName = conf.ShortName
	m.StartURL = conf.StartURL
	m.Display = conf.Display
	m.Description = conf.Description
	m.Dir = conf.Dir
	m.Orientation = conf.Orientation
	m.Scope = conf.Scope
	m.ThemeColor = conf.ThemeColor
	m.Background = conf.Background

	m.Icons = make([]*icon, len(conf.Icons))
	for index, img := range conf.Icons {
		m.Icons[index] = &icon{
			Src:   img.URL,
			Sizes: img.Sizes,
			Type:  img.Type,
		}
	}
}

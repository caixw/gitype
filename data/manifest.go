// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

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

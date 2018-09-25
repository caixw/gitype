// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/caixw/gitype/data/loader"
	"github.com/caixw/gitype/data/sw"
	"github.com/caixw/gitype/vars"
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

func (d *Data) buildManifest(conf *loader.Config) error {
	if conf.PWA == nil { // 不需要生成 pwa
		return nil
	}

	if conf.PWA.Manifest == nil {
		return nil
	}

	m := &Manifest{}
	m.fromLoader(conf.PWA.Manifest)

	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}

	d.Manifest = &Feed{
		URL:     conf.PWA.Manifest.URL,
		Type:    conf.PWA.Manifest.Type,
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

func (d *Data) buildSW(conf *loader.Config) error {
	sw := sw.New()

	// 首页、archives.html 和 tags.html
	ver := "gitype-" + strconv.FormatInt(d.Created.Unix(), 10)
	sw.Add(ver, "/", vars.TagsURL(), vars.ArchivesURL())

	for _, post := range d.Posts {
		ver = "post-" + strconv.FormatInt(post.Modified.Unix(), 10)
		sw.Add(ver, post.Permalink)
	}

	for _, tag := range d.Tags {
		ver = "tag-" + strconv.FormatInt(tag.Modified.Unix(), 10)
		sw.Add(ver, tag.Permalink)
	}

	// 主题提供的缓存内容
	ver = "theme-" + d.Theme.ID + "-" + d.Theme.Version
	for _, url := range d.Theme.Assets {
		if url == "" {
			continue
		}

		fmt.Println("url:", url)
		if !strings.HasPrefix(url, "https://") {
			url = themeURL(url)
		}
		sw.Add(ver, url)
	}

	d.ServiceWorker = sw.Bytes()
	d.ServiceWorkerPath = conf.PWA.ServiceWorker

	return nil
}

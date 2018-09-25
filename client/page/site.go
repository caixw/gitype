// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package page

import (
	"runtime"
	"time"

	"github.com/issue9/web"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/vars"
)

// Site 页面的附加信息，除非重新加载数据，否则内容不会变。
type Site struct {
	AppName    string // 程序名称
	AppURL     string // 程序官网
	AppVersion string // 当前程序的版本号
	GoVersion  string // 编译的 Go 版本号
	Theme      *data.Theme

	SiteName      string     // 网站名称
	Subtitle      string     // 网站副标题
	URL           string     // 网站地址，若是一个子目录，则需要包含该子目录
	Icon          *data.Icon // 网站图标
	Language      string     // 页面语言
	PostSize      int        // 总文章数量
	Beian         string     // 备案号
	Uptime        time.Time  // 上线时间
	LastUpdated   time.Time  // 最后更新时间
	RSS           *data.Link // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom          *data.Link
	Opensearch    *data.Link
	Manifest      *data.Link
	ServiceWorker string       // 指向 service worker 的 js 文件
	Tags          []*data.Tag  // 标签列表
	Series        []*data.Tag  // 专题列表
	Links         []*data.Link // 友情链接
	Menus         []*data.Link // 导航菜单
}

// NewSite 声明 Site 实例
func NewSite(d *data.Data) *Site {
	site := &Site{
		AppName:    vars.Name,
		AppURL:     vars.URL,
		AppVersion: vars.Version(),
		GoVersion:  runtime.Version(),
		Theme:      d.Theme,

		SiteName:      d.SiteName,
		Subtitle:      d.Subtitle,
		URL:           web.URL(""),
		Icon:          d.Icon,
		Language:      d.LanguageTag.String(),
		PostSize:      len(d.Posts),
		Beian:         d.Beian,
		Uptime:        d.Uptime,
		LastUpdated:   d.Created,
		ServiceWorker: d.ServiceWorkerPath,
		Tags:          d.Tags,
		Series:        d.Series,
		Links:         d.Links,
		Menus:         d.Menus,
	}

	if d.RSS != nil {
		site.RSS = &data.Link{
			Title: d.RSS.Title,
			URL:   d.RSS.URL,
			Type:  d.RSS.Type,
		}
	}

	if d.Atom != nil {
		site.Atom = &data.Link{
			Title: d.Atom.Title,
			URL:   d.Atom.URL,
			Type:  d.Atom.Type,
		}
	}

	if d.Opensearch != nil {
		site.Opensearch = &data.Link{
			Title: d.Opensearch.Title,
			URL:   d.Opensearch.URL,
			Type:  d.Opensearch.Type,
		}
	}

	if d.Manifest != nil {
		site.Manifest = &data.Link{
			Title: d.Manifest.Title,
			URL:   d.Manifest.URL,
			Type:  d.Manifest.Type,
		}
	}

	return site
}

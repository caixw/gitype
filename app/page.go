// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
	"github.com/issue9/logs"
)

// 用于描述一个页面的所有无素
type page struct {
	Title       string       // 文章标题，可以为空
	SiteName    string       // 网站名称
	Subtitle    string       // 副标题
	URL         string       // 网站主域名
	Root        string       // 网站的根目录
	Canonical   string       // 当前页的唯一链接
	Keywords    string       // meta.keywords的值
	Description string       // meta.description的值
	AppVersion  string       // 当前程序的版本号
	GoVersion   string       // 编译的go版本号
	PostSize    int          // 总文章数量
	Beian       string       // 备案号
	Uptime      int64        // 上线时间
	LastUpdated int64        // 最后更新时间
	RSS         *data.Link   // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom        *data.Link   // Atom
	PrevPage    *data.Link   // 前一页
	NextPage    *data.Link   // 下一页
	Tags        []*data.Tag  // 标签列表
	Links       []*data.Link // 友情链接
	Tag         *data.Tag    // 标签详细页面，非标签详细页，则为空
	Menus       []*data.Link // 菜单
	Posts       []*data.Post // 文章列表，文章列表页用到。
	Post        *data.Post   // 文章详细内容，单文章页面用到。

	app *app
}

func (a *app) newPage() *page {
	conf := a.data.Config

	page := &page{
		SiteName:    conf.Title,
		Subtitle:    conf.Subtitle,
		URL:         conf.URL,
		Root:        a.data.URLS.Root,
		Canonical:   conf.URL,
		Keywords:    conf.Keywords,
		Description: conf.Description,
		AppVersion:  vars.Version,
		GoVersion:   runtime.Version(),
		PostSize:    len(a.data.Posts),
		Beian:       conf.Beian,
		Uptime:      conf.Uptime,
		LastUpdated: a.updated,
		Tags:        a.data.Tags,
		Links:       a.data.Links,
		Menus:       conf.Menus,
		app:         a,
	}
	if conf.RSS != nil {
		page.RSS = &data.Link{Title: conf.RSS.Title, URL: conf.RSS.URL}
	}

	if conf.Atom != nil {
		page.Atom = &data.Link{Title: conf.Atom.Title, URL: conf.Atom.URL}
	}

	return page
}

// 输出当前内容到指定模板
func (p *page) render(w http.ResponseWriter, r *http.Request, name string, headers map[string]string) {
	for key, val := range headers {
		w.Header().Set(key, val)
	}

	err := p.app.data.Template.ExecuteTemplate(w, name, p)
	if err != nil {
		logs.Error("page.render:", err)
		p.renderStatusCode(w, r, http.StatusInternalServerError)
		return
	}
}

// 输出一个特定状态码下的错误页面。若该页面模板不存在，则panic。
// 只查找当前主题目录下的相关文件。
// 只对状态码大于等于400的起作用。
func (p *page) renderStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	if code < 400 {
		return
	}

	filename := strconv.Itoa(code) + ".html"
	path := filepath.Join(p.app.path.DataThemes, p.app.data.Config.Theme, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(code)
	w.Write(data)
}

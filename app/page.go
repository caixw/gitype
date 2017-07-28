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
	"github.com/issue9/utils"
)

const (
	contentTypeKey  = "Content-Type"
	contentTypeHTML = "text/html"
)

// 用于描述一个页面的所有无素
type page struct {
	a *app

	AppVersion  string       // 当前程序的版本号
	GoVersion   string       // 编译的 Go 版本号
	Title       string       // 文章标题，可以为空
	SiteName    string       // 网站名称
	Subtitle    string       // 副标题
	URL         string       // 网站主域名
	Language    string       // 页面语言
	Root        string       // 网站的根目录
	Canonical   string       // 当前页的唯一链接
	Keywords    string       // meta.keywords 的值
	Q           string       // 搜索关键字
	Description string       // meta.description 的值
	PostSize    int          // 总文章数量
	Beian       string       // 备案号
	Uptime      int64        // 上线时间
	LastUpdated int64        // 最后更新时间
	RSS         *data.Link   // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom        *data.Link   // Atom
	Opensearch  *data.Link   // Opensearch
	PrevPage    *data.Link   // 前一页
	NextPage    *data.Link   // 下一页
	Tags        []*data.Tag  // 标签列表
	Links       []*data.Link // 友情链接
	Tag         *data.Tag    // 标签详细页面，非标签详细页，则为空
	Menus       []*data.Link // 菜单
	Posts       []*data.Post // 文章列表，仅标签详情页和搜索页用到。
	Post        *data.Post   // 文章详细内容，仅文章页面用到。
}

func (a *app) newPage() *page {
	conf := a.buf.Data.Config

	page := &page{
		a:          a,
		AppVersion: vars.Version(),
		GoVersion:  runtime.Version(),

		SiteName:    conf.Title,
		Subtitle:    conf.Subtitle,
		Language:    conf.Language,
		URL:         conf.URL,
		Root:        "/",
		Canonical:   conf.URL,
		Keywords:    conf.Keywords,
		Description: conf.Description,
		PostSize:    len(a.buf.Data.Posts),
		Beian:       conf.Beian,
		Uptime:      conf.Uptime,
		LastUpdated: a.buf.Updated,
		Tags:        a.buf.Data.Tags,
		Links:       a.buf.Data.Links,
		Menus:       conf.Menus,
	}

	if conf.RSS != nil {
		page.RSS = &data.Link{Title: conf.RSS.Title, URL: conf.RSS.URL}
	}

	if conf.Atom != nil {
		page.Atom = &data.Link{Title: conf.Atom.Title, URL: conf.Atom.URL}
	}

	if conf.Opensearch != nil {
		page.Opensearch = &data.Link{Title: conf.Opensearch.Title, URL: conf.Opensearch.URL}
	}

	return page
}

// 输出当前内容到指定模板
func (p *page) render(w http.ResponseWriter, name string, headers map[string]string) {
	if len(headers) == 0 {
		w.Header().Set(contentTypeKey, contentTypeHTML)
	} else {
		if _, exists := headers[contentTypeKey]; !exists {
			headers[contentTypeKey] = contentTypeHTML
		}

		for key, val := range headers {
			w.Header().Set(key, val)
		}
	}

	err := p.a.buf.Template.ExecuteTemplate(w, name, p)
	if err != nil {
		logs.Error(err)
		p.a.renderError(w, http.StatusInternalServerError)
		return
	}
}

// 输出一个特定状态码下的错误页面。
// 若该页面模板不存在，则输出状态码对应的文本内容。
// 只查找当前主题目录下的相关文件。
// 只对状态码大于等于 400 的起作用。
func (a *app) renderError(w http.ResponseWriter, code int) {
	if code < 400 {
		return
	}
	logs.Debug("输出非正常状态码：", code)

	// 根据情况输出内容，若不存在模板，则直接输出最简单的状态码对应的文本。
	filename := strconv.Itoa(code) + ".html"
	path := filepath.Join(a.path.ThemesDir, a.buf.Data.Config.Theme, filename)
	if !utils.FileExists(path) {
		logs.Errorf("模板文件[%s]不存在\n", path)
		http.Error(w, http.StatusText(code), code)
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Errorf("读取模板文件[%v]时出现以下错误[%v]\n", path, err)
		http.Error(w, http.StatusText(code), code)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeHTML)
	w.WriteHeader(code)
	w.Write(data)
}

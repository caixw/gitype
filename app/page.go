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

// 定义页面的类型
const (
	typeIndex  = "index"
	typePosts  = "posts"
	typePost   = "post"
	typeTags   = "tags"
	typeTag    = "tag"
	typeSearch = "search"
	typeLinks  = "links"
)

const (
	contentTypeKey        = "Content-Type"
	contentTypeHTML       = "text/html"
	contentTypeXML        = "application/xml"
	contentTypeAtom       = "application/atom+xml"
	contentTypeRSS        = "application/rss+xml"
	contentTypeOpensearch = "application/opensearchdescription+xml"
)

// 生成一个带编码的 content-type 报头内容
func buildContentTypeContent(mime string) string {
	return mime + ";charset=utf-8"
}

func setContentType(w http.ResponseWriter, mime string) {
	w.Header().Set(contentTypeKey, buildContentTypeContent(mime))
}

// 用于描述一个页面的所有无素
type page struct {
	a    *app
	Info *info

	Title       string     // 文章标题，可以为空
	Subtitle    string     // 副标题
	Canonical   string     // 当前页的唯一链接
	Keywords    string     // meta.keywords 的值
	Q           string     // 搜索关键字
	Description string     // meta.description 的值
	PrevPage    *data.Link // 前一页
	NextPage    *data.Link // 下一页
	Type        string     // 当前页面类型

	// 以下内容，仅在对应的页面才会有内容
	Tag   *data.Tag    // 标签详细页面，非标签详细页，则为空
	Posts []*data.Post // 文章列表，仅标签详情页和搜索页用到。
	Post  *data.Post   // 文章详细内容，仅文章页面用到。
}

// 页面的附加信息，除非重新加载数据，否则内容不会变。
type info struct {
	AppName    string // 程序名称
	AppURL     string // 程序官网
	AppVersion string // 当前程序的版本号
	GoVersion  string // 编译的 Go 版本号

	ThemeName   string       // 主题名称
	ThemeURL    string       // 主题官网
	ThemeAuthor *data.Author // 主题的作者

	SiteName    string       // 网站名称
	URL         string       // 网站地址，若是一个子目录，则需要包含该子目录
	Icon        *data.Icon   // 网站图标
	Language    string       // 页面语言
	PostSize    int          // 总文章数量
	Beian       string       // 备案号
	Uptime      int64        // 上线时间
	LastUpdated int64        // 最后更新时间
	RSS         *data.Link   // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom        *data.Link   // Atom
	Opensearch  *data.Link   // Opensearch
	Tags        []*data.Tag  // 标签列表
	Links       []*data.Link // 友情链接
	Menus       []*data.Link // 菜单
}

func (a *app) newInfo() *info {
	conf := a.buf.Data.Config

	info := &info{
		AppName:    vars.AppName,
		AppURL:     vars.URL,
		AppVersion: vars.Version(),
		GoVersion:  runtime.Version(),

		ThemeName:   a.buf.Data.Theme.Name,
		ThemeURL:    a.buf.Data.Theme.URL,
		ThemeAuthor: a.buf.Data.Theme.Author,

		SiteName:    conf.Title,
		URL:         conf.URL,
		Icon:        conf.Icon,
		Language:    conf.Language,
		PostSize:    len(a.buf.Data.Posts),
		Beian:       conf.Beian,
		Uptime:      conf.Uptime,
		LastUpdated: a.buf.Updated,
		Tags:        a.buf.Data.Tags,
		Links:       a.buf.Data.Links,
		Menus:       conf.Menus,
	}

	if conf.RSS != nil {
		info.RSS = &data.Link{Title: conf.RSS.Title, URL: conf.RSS.URL}
	}

	if conf.Atom != nil {
		info.Atom = &data.Link{Title: conf.Atom.Title, URL: conf.Atom.URL}
	}

	if conf.Opensearch != nil {
		info.Opensearch = &data.Link{Title: conf.Opensearch.Title, URL: conf.Opensearch.URL}
	}

	return info
}

func (a *app) page(typ string) *page {
	conf := a.buf.Data.Config

	return &page{
		a:           a,
		Info:        a.info,
		Subtitle:    conf.Subtitle,
		Keywords:    conf.Keywords,
		Description: conf.Description,
		Type:        typ,
	}
}

func (p *page) nextPage(url, text string) {
	if len(text) == 0 {
		text = "下一页"
	}

	p.NextPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "next",
	}
}

func (p *page) prevPage(url, text string) {
	if len(text) == 0 {
		text = "上一页"
	}

	p.PrevPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "prev",
	}
}

// 输出当前内容到指定模板
func (p *page) render(w http.ResponseWriter, name string, headers map[string]string) {
	if len(headers) == 0 {
		setContentType(w, contentTypeHTML)
	} else {
		if _, exists := headers[contentTypeKey]; !exists {
			headers[contentTypeKey] = buildContentTypeContent(contentTypeHTML)
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
		statusError(w, code)
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Errorf("读取模板文件[%s]时出现以下错误[%v]\n", path, err)
		statusError(w, code)
		return
	}

	setContentType(w, contentTypeHTML)
	w.WriteHeader(code)
	w.Write(data)
}

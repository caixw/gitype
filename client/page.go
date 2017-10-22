// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/caixw/gitype/data"
	"github.com/caixw/gitype/helper"
	"github.com/caixw/gitype/vars"
	"github.com/issue9/logs"
	"github.com/issue9/utils"
)

const contentTypeKey = "Content-Type"

// 生成一个带编码的 content-type 报头内容
func buildContentTypeContent(mime string) string {
	return mime + ";charset=utf-8"
}

// 设置页面的编码，若已经存在，则不会受影响。
// 要强制指定，请直接使用 w.Header().Set()
func setContentType(w http.ResponseWriter, mime string) {
	h := w.Header()
	if len(h.Get(contentTypeKey)) == 0 {
		h.Set(contentTypeKey, buildContentTypeContent(mime))
	}
}

// 用于描述一个页面的所有无素
type page struct {
	client   *Client
	Info     *info
	response http.ResponseWriter
	request  *http.Request

	Title       string       // 文章标题
	Subtitle    string       // 副标题
	Canonical   string       // 当前页的唯一链接
	Keywords    string       // meta.keywords 的值
	Description string       // meta.description 的值
	PrevPage    *data.Link   // 前一页
	NextPage    *data.Link   // 下一页
	Type        string       // 当前页面类型
	Author      *data.Author // 作者
	License     *data.Link   // 当前页的版本信息，可以为空

	// 以下内容，仅在对应的页面才会有内容
	Q        string          // 搜索关键字
	Tag      *data.Tag       // 标签详细页面，非标签详细页，则为空
	Posts    []*data.Post    // 文章列表，仅标签详情页和搜索页用到。
	Post     *data.Post      // 文章详细内容，仅文章页面用到。
	Archives []*data.Archive // 归档
}

// 页面的附加信息，除非重新加载数据，否则内容不会变。
type info struct {
	AppName    string // 程序名称
	AppURL     string // 程序官网
	AppVersion string // 当前程序的版本号
	GoVersion  string // 编译的 Go 版本号
	Theme      *data.Theme

	SiteName    string     // 网站名称
	URL         string     // 网站地址，若是一个子目录，则需要包含该子目录
	Icon        *data.Icon // 网站图标
	Language    string     // 页面语言
	PostSize    int        // 总文章数量
	Beian       string     // 备案号
	Uptime      time.Time  // 上线时间
	LastUpdated time.Time  // 最后更新时间
	RSS         *data.Link // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom        *data.Link
	Opensearch  *data.Link
	Tags        []*data.Tag  // 标签列表
	Series      []*data.Tag  // 专题列表
	Links       []*data.Link // 友情链接
	Menus       []*data.Link // 导航菜单
}

func (client *Client) newInfo() *info {
	d := client.data

	info := &info{
		AppName:    vars.Name,
		AppURL:     vars.URL,
		AppVersion: vars.Version(),
		GoVersion:  runtime.Version(),
		Theme:      d.Theme,

		SiteName:    d.SiteName,
		URL:         d.URL,
		Icon:        d.Icon,
		Language:    d.Language,
		PostSize:    len(d.Posts),
		Beian:       d.Beian,
		Uptime:      d.Uptime,
		LastUpdated: d.Created,
		Tags:        d.Tags,
		Series:      d.Series,
		Links:       d.Links,
		Menus:       d.Menus,
	}

	if d.RSS != nil {
		info.RSS = &data.Link{
			Title: d.RSS.Title,
			URL:   d.RSS.URL,
			Type:  d.RSS.Type,
		}
	}

	if d.Atom != nil {
		info.Atom = &data.Link{
			Title: d.Atom.Title,
			URL:   d.Atom.URL,
			Type:  d.Atom.Type,
		}
	}

	if d.Opensearch != nil {
		info.Opensearch = &data.Link{
			Title: d.Opensearch.Title,
			URL:   d.Opensearch.URL,
			Type:  d.Opensearch.Type,
		}
	}

	return info
}

func (client *Client) page(typ string, w http.ResponseWriter, r *http.Request) *page {
	d := client.data

	return &page{
		client:   client,
		Info:     client.info,
		response: w,
		request:  r,

		Subtitle: d.Subtitle,
		Type:     typ,
		Author:   d.Author,
		License:  d.License,
	}
}

func (p *page) nextPage(url, text string) {
	p.NextPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "next",
	}
}

func (p *page) prevPage(url, text string) {
	p.PrevPage = &data.Link{
		Text: text,
		URL:  url,
		Rel:  "prev",
	}
}

// 输出当前内容到指定模板
func (p *page) render(name string) {
	setContentType(p.response, p.client.data.Type)

	err := p.client.data.Theme.Template.ExecuteTemplate(p.response, name, p)
	if err != nil {
		logs.Error(err)
		p.client.renderError(p.response, p.request, http.StatusInternalServerError)
		return
	}
}

// 输出一个特定状态码下的错误页面。
// 若该页面模板不存在，则输出状态码对应的文本内容。
// 只查找当前主题目录下的相关文件。
// 只对状态码大于等于 400 的起作用。
func (client *Client) renderError(w http.ResponseWriter, r *http.Request, code int) {
	if code < 400 {
		return
	}
	logs.Debug("输出非正常状态码：", code)

	// 根据情况输出内容，若不存在模板，则直接输出最简单的状态码对应的文本。
	filename := strconv.Itoa(code) + vars.TemplateExtension
	path := filepath.Join(client.path.ThemesDir, client.data.Theme.ID, filename)
	if !utils.FileExists(path) {
		logs.Debugf("模板文件 %s 不存在\n", path)
		helper.StatusError(w, code)
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Errorf("读取模板文件 %s 时出现以下错误: %v\n", path, err)
		helper.StatusError(w, code)
		return
	}

	setContentType(w, client.data.Type)
	w.WriteHeader(code)
	w.Write(data)
}

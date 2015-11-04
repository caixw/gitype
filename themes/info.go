// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"runtime"
	"strconv"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm/fetch"
)

type Anchor struct {
	Link  string
	Title string
}

// 页面的基本信息
type Info struct {
	Title       string   // 网页的title值
	SiteURL     string   // 网站地址
	SiteName    string   // 网站名称
	SecondTitle string   // 副标题
	Canonical   string   // 当前页的唯一链接
	RSSURL      string   // RSS
	AtomURL     string   // Atom
	Keywords    string   // meta.keywords的值
	Description string   // meta.description的值
	AppVersion  string   // 当前程序的版本号
	GoVersion   string   // 编译的go版本号
	CurrentPage int      // 当前页码
	PostSize    int      // 总文章数量
	CommentSize int      // 总评论数量
	Tags        []*Tag   // 标签列表
	Tops        []*Post  // 最新评论的10条内容
	Hots        []*Post  // 评论最多的10条内容
	Menus       []Anchor // 菜单
}

func getInfo() (*Info, error) {
	p := &Info{
		SiteURL:     opt.SiteURL,
		SiteName:    opt.SiteName,
		SecondTitle: opt.SecondTitle,
		Keywords:    opt.Keywords,
		Description: opt.Description,
		AppVersion:  core.Version,
		GoVersion:   runtime.Version(),
		Menus: []Anchor{ // TODO 添加到options配置中
			{Link: opt.SiteURL, Title: "首页"},
			{Link: core.PostURL(opt, "about"), Title: "关于"},
			{Link: core.TagsURL(opt), Title: "标签"},
		},
	}

	var err error
	sql := "SELECT COUNT(*) as cnt FROM #posts WHERE {state}=?" // TODO 预编译成stmt
	if p.PostSize, err = getSize(sql, models.PostStatePublished); err != nil {
		return nil, err
	}

	sql = "SELECT COUNT(*) as cnt FROM #comments WHERE {state}=?"
	if p.CommentSize, err = getSize(sql, models.CommentStateApproved); err != nil {
		return nil, err
	}

	if p.Tags, err = getTags(); err != nil {
		return nil, err
	}

	if p.Tops, err = getTops(); err != nil {
		return nil, err
	}

	return p, nil
}

func getSize(sql string, args ...interface{}) (int, error) {
	rows, err := db.Query(true, sql, args...)
	if err != nil {
		return 0, err
	}
	cnts, err := fetch.ColumnString(true, "cnt", rows)
	rows.Close()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cnts[0])
}

func getTops() ([]*Post, error) {
	sql := `SELECT c.{content} AS Content, p.{title} AS Tilte, p.{name} AS Name, p.{id} AS ID
	FROM #comments AS c
	LEFT JOIN #posts AS p ON c.{postID}=p.{id}
	WHERE c.{state}=?
	ORDER BY c.{id} DESC
	LIMIT ? `
	rows, err := db.Query(true, sql, models.CommentStateApproved, opt.SidebarSize)
	if err != nil {
		return nil, err
	}

	topics := make([]*Post, 0, opt.SidebarSize)
	_, err = fetch.Obj(&topics, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func getTags() ([]*Tag, error) {
	sql := "SELECT  {title} AS Title, {name} AS Name FROM #tags"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0, opt.SidebarSize)
	_, err = fetch.Obj(&tags, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	return tags, nil
}

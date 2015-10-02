// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

// State 表示一些当前的博客状态。
type State struct {
	Posts          int // 文章的数量
	DraftPosts     int // 草稿的数量
	PublishedPosts int // 已发布文章的数量

	Comments         int // 评论数量
	WaitingComments  int // 待审评论数量
	SpamComments     int // 垃圾评论数量
	ApprovedComments int // 已审评论数量

	LastLogin  int    // 最后次登录时间
	LastPosted int    // 最后次发表文章的时间
	LastIP     string // 最后次登录的IP
	LastAgent  string // 最后次登录的浏览器相关资料

	ScreenName string // 用户的当前昵称

	CurrentTheme string // 当前的主题
	Themes       int    // 的有的主题数量
}

// 重新加载state数据
/*func LoadState(db *orm.DB) (*State, error) {
	state := &State{}
	var err error

	sql := "select count(*) as cnt,{state} from #post"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data, err := fetch.MapString(false, rows)
	if err != nil {
		return nil, err
	}
	// TODO

	return state, nil
}*/

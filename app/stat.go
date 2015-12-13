// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"github.com/caixw/typing/models"
	"github.com/issue9/orm"
)

// 一些临时性的统计数据，在程序启动时初始化，关闭之后也不会被保存到数据库。
type Stat struct {
	CommentsSize         int           // 评论数
	WaitingCommentsSize  int           // 待评论数量
	ApprovedCommentsSize int           // 待评论数量
	SpamCommentsSize     int           // 垃圾论数量
	PostsSize            int           // 文章数量
	PublishedPostsSize   int           // 已发表文章数量
	DraftPostsSize       int           // 草稿数量
	Posts                map[int64]int // 文章对应的评论数量
	Tags                 map[int64]int // 标签对应的文章数量
}

func loadStat(db *orm.DB) (*Stat, error) {
	stat := &Stat{}
	var err error

	/* comments */
	o := &models.Comment{State: models.CommentStateSpam}
	stat.SpamCommentsSize, err = db.Count(o)
	if err != nil {
		return nil, err
	}

	o.State = models.CommentStateWaiting
	stat.WaitingCommentsSize, err = db.Count(o)
	if err != nil {
		return nil, err
	}

	o.State = models.CommentStateApproved
	stat.ApprovedCommentsSize, err = db.Count(o)
	if err != nil {
		return nil, err
	}

	stat.CommentsSize = stat.SpamCommentsSize + stat.WaitingCommentsSize + stat.ApprovedCommentsSize

	/* posts */
	p := &models.Post{State: models.PostStateDraft}
	stat.DraftPostsSize, err = db.Count(p)
	if err != nil {
		return nil, err
	}

	p.State = models.PostStatePublished
	stat.PublishedPostsSize, err = db.Count(p)
	if err != nil {
		return nil, err
	}

	stat.PostsSize = stat.PublishedPostsSize + stat.DraftPostsSize

	return stat, nil
}

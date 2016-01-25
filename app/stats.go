// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import "github.com/caixw/typing/models"

// 一些临时性的统计数据，在程序启动时初始化，关闭之后也不会被保存到数据库。
type Stats struct {
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

func GetStats() *Stats {
	return stats
}

// 从数据库初始化数据
func loadStats() (*Stats, error) {
	stats := &Stats{}

	if err := stats.ReBuild(); err != nil {
		return nil, err
	}

	return stats, nil
}

// 重新构建数据
func (s *Stats) ReBuild() error {
	/* 统计评论数量 */
	if err := s.UpdateCommentsSize(); err != nil {
		return err
	}

	/* 统计文章数量 */
	if err := s.UpdatePostsSize(); err != nil {
		return err
	}

	// posts
	s.Posts = make(map[int64]int, s.PostsSize)
	s.Tags = make(map[int64]int, 100)

	return nil
}

//更新文章评论数量
func (s *Stats) UpdatePostsSize() (err error) {
	p := &models.Post{State: models.PostStateDraft}
	s.DraftPostsSize, err = db.Count(p)
	if err != nil {
		return err
	}

	p.State = models.PostStatePublished
	s.PublishedPostsSize, err = db.Count(p)
	if err != nil {
		return err
	}

	s.PostsSize = s.PublishedPostsSize + s.DraftPostsSize

	return nil
}

// 更新评论数据
func (s *Stats) UpdateCommentsSize() (err error) {
	o := &models.Comment{State: models.CommentStateSpam}
	s.SpamCommentsSize, err = db.Count(o)
	if err != nil {
		return err
	}

	o.State = models.CommentStateWaiting
	s.WaitingCommentsSize, err = db.Count(o)
	if err != nil {
		return err
	}

	o.State = models.CommentStateApproved
	s.ApprovedCommentsSize, err = db.Count(o)
	if err != nil {
		return err
	}

	s.CommentsSize = s.SpamCommentsSize + s.WaitingCommentsSize + s.ApprovedCommentsSize

	return nil
}

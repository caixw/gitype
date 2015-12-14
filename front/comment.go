// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"strconv"

	"github.com/caixw/typing/app"
)

type Comment struct {
	PostID   int64
	PostName string

	ID          int64
	Created     int64
	IP          string
	Agent       string
	Content     string
	IsAdmin     bool
	AuthorName  string
	AuthorURL   string
	AuthorEmail string
}

// 文章的固定链接，相对URL，若要绝对URL，请使用opt.URL()进行封装。
func (c *Comment) Permalink() string {
	if len(c.PostName) > 0 {
		return opt.CommentURL(c.PostName, c.ID)
	}
	return opt.CommentURL(strconv.FormatInt(c.PostID, 10), c.ID)
}

func (c *Comment) Fragment() string {
	return app.CommentFragment(c.ID)
}

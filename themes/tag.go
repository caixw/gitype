// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import "github.com/caixw/typing/core"

type Tag struct {
	ID          int64
	Name        string
	Title       string
	Description string
	Count       int // 关联的文章数量
}

func (t *Tag) Permalink() string {
	return core.TagURL(t.Name, 1)
}

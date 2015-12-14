// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
)

type Tag struct {
	ID          int64
	Name        string
	Title       string
	Description string
	//Count       int // 关联的文章数量
}

func (t *Tag) Count() int {
	if cnt, found := stat.Tags[t.ID]; found {
		return cnt
	}

	r := &models.Relationship{TagID: t.ID}
	cnt, err := db.Count(r)
	if err != nil {
		logs.Error("themes.Tag.Count:", err)
	}
	stat.Tags[t.ID] = cnt
	return cnt
}

func (t *Tag) Permalink() string {
	return opt.TagURL(t.Name, 1)
}

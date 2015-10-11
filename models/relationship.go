// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package models

type Relationship struct {
	PostID int64 `orm:"name(postID);pk"`
	TagID  int64 `orm:"name(tagID);pk"`
}

func (r *Relationship) Meta() string {
	return `name(relationships)`
}

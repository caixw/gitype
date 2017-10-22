// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"strconv"
	"time"
)

// Etag 根据一个时间，生成一段 Etag 字符串
func Etag(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

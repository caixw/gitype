// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"strconv"

	"github.com/caixw/typing/core"
)

// 为post生成一条唯一URL
func PostURL(opt *core.Options, p *Post) string {
	if len(p.Name) > 0 {
		return opt.SiteURL + "/posts/" + p.Name + opt.Suffix
	}

	return opt.SiteURL + "/posts/" + strconv.FormatInt(p.ID, 10) + opt.Suffix
}

func TagURL(opt *core.Options, name string) string {
	return opt.SiteURL + "/tags/" + name
}

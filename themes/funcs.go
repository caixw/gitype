// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"html/template"
	"time"
)

func dateFormat(t int64) interface{} {
	return time.Unix(t, 0).Format(opt.DateFormat)
}

func htmlEscaped(html string) interface{} {
	return template.HTML(html)
}

var funcMap = template.FuncMap{
	"html": htmlEscaped,
	"date": dateFormat,
}

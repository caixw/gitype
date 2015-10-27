// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import "html/template"

func htmlEscaped(x string) interface{} { return template.HTML(x) }

var funcs = template.FuncMap{"html": htmlEscaped}

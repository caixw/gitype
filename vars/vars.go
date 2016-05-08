// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// vars 包定义一些全局变量、常量。
package vars

const (
	// 版本号
	Version = "0.2.33.160508"

	// 媒体文件的地址前缀，相对于URLS.Root地址。
	// 之所以不配置在data/URLS中，是因为如果修改该值，
	// 会造成所有文章中对这些文件的引用都要修改，造成不必要的麻烦。
	MediaURL = "/media"

	// 客户配置文件中所使用的的时间格式。
	// 所有的时间字符串，都将使用此格式去解析。
	DateFormat = "2006-01-02T15:04:05-0700"
)

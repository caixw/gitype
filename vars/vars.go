// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 定义一些程序通用的数据。
package vars

// 版本号
const Version = "0.1.11.20160305"

const (
	// 媒体文件的地址前缀，之所以不配置在data/URLS中，是因为
	// 如果修改该值，会造成所有文章中对这些文章的引用都要修改，
	// 造成不必要的麻烦。
	// 相对于URLS.Root地址。
	MediaURL = "/media"

	AdminURL      = "/admin"
	AdminPassword = "123"
)

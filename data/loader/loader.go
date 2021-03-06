// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package loader 用于加载原始的数据内容
//
// loader 只负责加载数据，而对数据的处理则由 data 包负责。
package loader

const (
	contentTypeAtom       = "application/atom+xml"
	contentTypeRSS        = "application/rss+xml"
	contentTypeOpensearch = "application/opensearchdescription+xml"
	contentTypeXML        = "application/xml"
	contentManifest       = "application/manifest+json"
	contentTypeHTML       = "text/html"
)

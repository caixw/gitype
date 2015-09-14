// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/issue9/logs"
)

func pageIndex(w http.ResponseWriter, r *http.Request) {
	if err := tpl.Execute(w, nil); err != nil {
		logs.Error("pageIndex:", err)
	}
}

func pageTags(w http.ResponseWriter, r *http.Request) {

}

func pageTag(w http.ResponseWriter, r *http.Request) {

}

func pageCats(w http.ResponseWriter, r *http.Request) {

}

func pageCat(w http.ResponseWriter, r *http.Request) {

}

func pagePosts(w http.ResponseWriter, r *http.Request) {

}

func pagePost(w http.ResponseWriter, r *http.Request) {

}

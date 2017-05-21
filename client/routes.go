// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"path/filepath"

	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
	"github.com/issue9/utils"
)

var ()

func (c *Client) removeRoutes() {
	for _, route := range c.routes {
		c.mux.Remove(route)
	}

	c.routes = nil
}

func (c *Client) initRoutes() {
	pattern := vars.Post + "/{slug}" + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getPost))

	pattern = vars.Posts + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getPosts))

	pattern = vars.Tag + "/{slug}" + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getTag))

	pattern = vars.Tags + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getTags))

	pattern = vars.Search + vars.Suffix
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getSearch))

	pattern = vars.Themes + "/*"
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getThemes))

	pattern = "/*"
	c.routes = append(c.routes, pattern)
	c.mux.GetFunc(pattern, c.pre(c.getRaws))
}

func (c *Client) getPost(w http.ResponseWriter, r *http.Request) {
	ps := mux.GetParams(r)
	if ps == nil {
		// TODO
	}
}

func (c *Client) getPosts(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getTag(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getTags(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getThemes(w http.ResponseWriter, r *http.Request) {
	//
}

func (c *Client) getSearch(w http.ResponseWriter, r *http.Request) {
	//
}

// 读取根下的文件
func (c *Client) getRaws(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		c.getPosts(w, r)
		return
	}

	root := http.Dir(c.path.RawsDir)
	if !utils.FileExists(filepath.Join(c.path.RawsDir, r.URL.Path)) {
		c.renderStatusCode(w, http.StatusNotFound)
		return
	}
	prefix := "/"
	http.StripPrefix(prefix, http.FileServer(root)).ServeHTTP(w, r)

}
